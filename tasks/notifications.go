package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/labstack/gommon/log"
	"strconv"
	"test-manager/monitoring"
	"test-manager/repos"
	"test-manager/services/alert_system"
	"test-manager/tasks/task_models"
	"test-manager/usecase_models"
	"text/template"
)

type NotificationTaskHandler struct {
	alertSystem alert_system.AlertHandler
	projectRepo repos.ProjectsRepository
}

func NewNotificationTaskHandler(
	alertSystem alert_system.AlertHandler,
	projectRepo repos.ProjectsRepository,
) *NotificationTaskHandler {
	return &NotificationTaskHandler{
		alertSystem: alertSystem,
		projectRepo: projectRepo,
	}
}

func (c *NotificationTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload task_models.NotificationsPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed on endpoint task payload: %v: %w", err, asynq.SkipRetry)
	}

	project, err := c.projectRepo.GetProject(ctx, payload.ProjectId)
	if err != nil {
		return err
	}

	var notifications usecase_models.Notifications
	err = json.Unmarshal(project.Notifications.JSON, &notifications)
	if err != nil {
		log.Error("can not unmarshal notification")
		return fmt.Errorf("json.Unmarshal failed on endpoint task notif json: %v: %w", err, asynq.SkipRetry)
	}

	subject := payload.PipelineName
	slackMessage := ""
	telegramMessage := ""
	emailMessage := ""
	emailKeyPairs := map[string]string{}
	templateKey := ""
	switch payload.State {
	case "up":
		subject = payload.PipelineName + " is UP"
		slackMessage = fmt.Sprintf("Hi %s,\n%s (%s) is up again\nresolved address:%s\n\nDatacenters envolved: %s\n\nFixed at %s",
			payload.Username, payload.PipelineName, payload.Type, payload.Address, payload.Datacenters, payload.Time)
		telegramMessage = fmt.Sprintf("Hi %s,\n%s (%s) is up again\nresolved address:%s\n\nDatacenters envolved: %s\n\nFixed at %s",
			payload.Username, payload.PipelineName, payload.Type, payload.Address, payload.Datacenters, payload.Time)
		//emailMessage, _ = ParseTemplate("resolved.html", TemplateKeys{
		//	Username:     payload.Username,
		//	Address:      payload.Address,
		//	PipelineName: payload.PipelineName,
		//	Datacenters:  payload.Datacenters,
		//	RootCause:    payload.RootCause,
		//	Time:         payload.Time,
		//})
		emailKeyPairs = map[string]string{
			"object_name":          payload.PipelineName,
			"address":              payload.Address,
			"root_cause":           payload.RootCause,
			"datacenters_resolved": payload.Datacenters,
			"incident_start_time":  payload.Time,
			"incident_end_time":    payload.Time,
			"incident_duration":    payload.Time,
		}
		templateKey = "resolved"
		emailMessage = telegramMessage
	case "down":
		subject = payload.PipelineName + " is DOWN"
		slackMessage = fmt.Sprintf("Hi %s,\n%s (%s) is DOWN\nchecked address:%s\n\nDatacenters envolved: %s\n\nIncident at %s",
			payload.Username, payload.PipelineName, payload.Type, payload.Address, payload.Datacenters, payload.Time)
		telegramMessage = fmt.Sprintf("Hi %s,\n%s (%s) is DOWN\nchecked address:%s\n\nDatacenters envolved: %s\n\nIncident at %s",
			payload.Username, payload.PipelineName, payload.Type, payload.Address, payload.Datacenters, payload.Time)
		//emailMessage, _ = ParseTemplate("detected.html", TemplateKeys{
		//	Username:     payload.Username,
		//	Address:      payload.Address,
		//	PipelineName: payload.PipelineName,
		//	Datacenters:  payload.Datacenters,
		//	RootCause:    payload.RootCause,
		//	Time:         payload.Time,
		//})
		emailKeyPairs = map[string]string{
			"object_name":         payload.PipelineName,
			"address":             payload.Address,
			"root_cause":          payload.RootCause,
			"datacenters_failed":  payload.Datacenters,
			"incident_start_time": payload.Time,
		}
		templateKey = "problem_detected"
		emailMessage = telegramMessage
	case "diff":
		subject = payload.PipelineName + " is DOWN update"
		slackMessage = fmt.Sprintf("Hi %s,\n%s (%s) is DOWN update\nchecked address:%s\n\nDatacenters fixed: %s\nNew Datacenters failed:%s\n\nUpdated at %s",
			payload.Username, payload.PipelineName, payload.Type, payload.Address, payload.ResolvedDatacenters, payload.FailedDatacenters, payload.Time)
		telegramMessage = fmt.Sprintf("Hi %s,\n%s (%s) is DOWN update\nchecked address:%s\n\nDatacenters fixed: %s\nNew Datacenters failed:%s\n\nUpdated at %s",
			payload.Username, payload.PipelineName, payload.Type, payload.Address, payload.ResolvedDatacenters, payload.FailedDatacenters, payload.Time)
		//emailMessage, _ = ParseTemplate("detected_update.html", TemplateKeys{
		//	Username:            payload.Username,
		//	Address:             payload.Address,
		//	PipelineName:        payload.PipelineName,
		//	Datacenters:         payload.Datacenters,
		//	FailedDatacenters:   payload.FailedDatacenters,
		//	ResolvedDatacenters: payload.ResolvedDatacenters,
		//	RootCause:           payload.RootCause,
		//	Time:                payload.Time,
		//})
		emailKeyPairs = map[string]string{
			"object_name":          payload.PipelineName,
			"address":              payload.Address,
			"root_cause":           payload.RootCause,
			"datacenters_resolved": payload.ResolvedDatacenters,
			"datacenters_failed":   payload.FailedDatacenters,
			"incident_start_time":  payload.Time,
		}
		templateKey = "problem_detected_update"
		emailMessage = telegramMessage
	}

	monitoring.NotificationTaskCounter.WithLabelValues(payload.State).Inc()

	if len(notifications.Slack) != 0 {
		err = c.alertSystem.SendAlert(ctx, alert_system.AlertRequest{
			AlertType:        "slack",
			UserId:           strconv.Itoa(payload.ProjectId),
			Targets:          notifications.Slack,
			Subject:          subject,
			Message:          slackMessage,
			IsTemplate:       true,
			Template:         "",
			TemplateKeyPairs: nil,
			AdditionalData:   nil,
		})
		if err != nil {
			log.Info("error on sending slack alert in executing rule: ", err)
		}
	}
	if len(notifications.Email) != 0 {
		err = c.alertSystem.SendAlert(ctx, alert_system.AlertRequest{
			AlertType:        "email",
			UserId:           strconv.Itoa(payload.ProjectId),
			Targets:          notifications.Email,
			Subject:          subject,
			Message:          emailMessage,
			IsTemplate:       true,
			Template:         templateKey,
			TemplateKeyPairs: emailKeyPairs,
			AdditionalData:   nil,
		})
		if err != nil {
			log.Info("error on sending email alert in executing rule: ", err)
		}
	}
	if len(notifications.Telegram) != 0 {
		err = c.alertSystem.SendAlert(ctx, alert_system.AlertRequest{
			AlertType: "telegram",
			Subject:   subject,
			Message:   telegramMessage,
			UserId:    strconv.Itoa(payload.ProjectId),
			Targets:   notifications.Telegram,
		})
		if err != nil {
			log.Info("error on sending telegram alert in executing rule: ", err)
		}
	}

	fmt.Println("success on processing notification task")
	return nil
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles("./templates/" + templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type TemplateKeys struct {
	Username            string
	Address             string
	PipelineName        string
	Datacenters         string
	FailedDatacenters   string
	ResolvedDatacenters string
	RootCause           string
	Time                string
}

var ProblemDetectedEmailTemplate = `<html>
    <head>
          <meta name=3D"viewport" =
content=3D"width=3Ddevice-width">
          <meta http-equiv=3D"Content-Typ=
e" content=3D"text/html; charset=3DUTF-8">
            <title>Monitor is  =
DOWN : {{.PipelineName}}
                    </title>
    </head>
    <body style=3D"font-family: 'Roboto', Arial, sans-serif; color: =
#131a26; background: #fefefe; line-height: 1.3em;">
      <div style=3D"font-family: 'Roboto', Arial, sans-serif; color: =
#131a26; background: #fefefe; box-sizing: border-box; line-height: 1.3em;">
        <div style=3D"background: #131a26; padding: 45px 15px;">
          <div style=3D"max-width: 600px; margin: 0 auto;">
            <table style=3D" width: 100%; font-size: 14px; color: #ffffff;"=
 width=3D"100%">
              <tbody style=3D"">
                <tr =
style=3D"">
                  <td width=3D"180" style=3D" font-size: =
36px;">
                    <a href=3D"https://uptimerobot.com/dashboard?=
utm_source=3DalertMessage&utm_medium=3Demail&utm_campaign=3Ddown-5x-free&ut=
m_content=3DheaderLogo#mainDashboard" style=3D" color: #3bd671 !=
important;"><img src=3D"https://cdn.mcauto-images-production.sendgrid.=
net/3e8054f26aace367/6a696f89-4c68-4736-83c3-c08df64e1a0d/549x79.png" =
alt=3D"UptimeRobot" width=3D"180" style=3D" max-width: 180px; display: =
inline-block;"></a>
                  </td>
                </tr>
              </tbody>
            </table>
            <div>
              <h1 style=3D"font-size: 36px; color: =
#ffffff; margin-top: 45px; margin-bottom: 10px; line-height: 28px;">
                    {{.PipelineName}} is <span style=3D"color:#df484a">down</span>.
                              </h1>
            </div>
          </div>
        </div>
        <div style=3D"max-width: 600px; margin: 30px auto 0 =
auto; padding: 0 10px;">

          <div style=3D"background: #ffffff; =
border-radius: 6px; box-shadow: 0 20px 40px 0 rgba(0,0,0,0.1); padding: =
25px; border: 1px solid #efefef; font-size: 14px; margin-bottom: 25px;">
            <div>

              <p style=3D"line-height: 20px;">Hello {{.Username}},</p>
                  <p style=3D"line-height: 20px;"><strong>We=
 just detected an incident on your monitor. Your service is currently down.=
</strong></p>
                  <p style=3D"line-height: 20px;">We will =
alert you when it&apos;s up again.</p>
                          </div>

            <div style=3D"padding: 20px; background-color: #f9f9f9; =
border-radius: 6px; margin-top: 25px; line-height: 16px;">

              <div style=3D"font-size: 12px; color: #687790;">Monitor =
name</div>
              <h2 style=3D"font-size: 14px; margin-bottom: 5px; =
margin-top: 3px;">{{.PipelineName}}</h2>
              <hr style=3D"border: 1px solid =
#dedede; border-bottom: 0; margin: 10px 0;">

              <div =
style=3D"color: #687790; font-size: 12px;">Checked URL</div>
              <h2 style=3D"font-size: 14px; margin-bottom: 5px; margin-top:=
 3px; line-height: 16px;"><code>
                <a href=3D"{{.Address}}" style=3D"color: #131a26 !important; text-decoration: none !=
important;">{{.Address}}</a></code></h2>
              <hr style=3D"border: 1px solid #dedede; border-bottom: 0; =
margin: 10px 0;">

<div style=3D"color: #687790; font-size: =
12px;">Datacenters Showed Problem</div>
<h2 style=3D"font-size: 14px; =
margin-bottom: 5px; margin-top: 3px;">
                  {{.Datacenters}}
              </h2>
              <hr style=3D"border: 1px solid #dedede; =
border-bottom: 0; margin: 10px 0;">


              <div style=3D"color: #687790; font-size: =
12px;">Root cause</div>
              <h2 style=3D"font-size: 14px; =
margin-bottom: 5px; margin-top: 3px;">
                  {{.RootCause}}
              </h2>
              <hr style=3D"border: 1px solid #dedede; =
border-bottom: 0; margin: 10px 0;">

              <div style=3D"font-size:=
 12px; color: #687790;">Incident started at</div>
              <h2 =
style=3D"font-size: 14px; margin-bottom: 5px; margin-top: 3px;">{{.Time}}</h2>
              <hr style=3D"border: 1px solid #dedede; =
border-bottom: 0; margin: 10px 0;">
            </div>
          </div>
        </div>
      </div>
  </body>
</html>`

var ProblemResolvedEmailTemplate = `<html>
    <head>
          <meta name=3D"viewport" =
content=3D"width=3Ddevice-width">
          <meta http-equiv=3D"Content-Typ=
e" content=3D"text/html; charset=3DUTF-8">
            <title>Monitor is  =
UP : {{.PipelineName}}
                    </title>
    </head>
    <body style=3D"font-family: 'Roboto', Arial, sans-serif; color: =
#131a26; background: #fefefe; line-height: 1.3em;">
      <div style=3D"font-family: 'Roboto', Arial, sans-serif; color: =
#131a26; background: #fefefe; box-sizing: border-box; line-height: 1.3em;">
        <div style=3D"background: #131a26; padding: 45px 15px;">
          <div style=3D"max-width: 600px; margin: 0 auto;">
            <table style=3D" width: 100%; font-size: 14px; color: #ffffff;"=
 width=3D"100%">
              <tbody style=3D"">
                <tr =
style=3D"">
                  <td width=3D"180" style=3D" font-size: =
36px;">
                    <a href=3D"https://uptimerobot.com/dashboard?=
utm_source=3DalertMessage&utm_medium=3Demail&utm_campaign=3Ddown-5x-free&ut=
m_content=3DheaderLogo#mainDashboard" style=3D" color: #3bd671 !=
important;"><img src=3D"https://cdn.mcauto-images-production.sendgrid.=
net/3e8054f26aace367/6a696f89-4c68-4736-83c3-c08df64e1a0d/549x79.png" =
alt=3D"UptimeRobot" width=3D"180" style=3D" max-width: 180px; display: =
inline-block;"></a>
                  </td>
                </tr>
              </tbody>
            </table>
            <div>
              <h1 style=3D"font-size: 36px; color: =
#ffffff; margin-top: 45px; margin-bottom: 10px; line-height: 28px;">
                    {{.PipelineName}} is <span style=3D"color:#df484a">down</span>.
                              </h1>
            </div>
          </div>
        </div>
        <div style=3D"max-width: 600px; margin: 30px auto 0 =
auto; padding: 0 10px;">

          <div style=3D"background: #ffffff; =
border-radius: 6px; box-shadow: 0 20px 40px 0 rgba(0,0,0,0.1); padding: =
25px; border: 1px solid #efefef; font-size: 14px; margin-bottom: 25px;">
            <div>

              <p style=3D"line-height: 20px;">Hello {{.Username}},</p>
                  <p style=3D"line-height: 20px;"><strong>We=
 just detected an incident on your monitor. Your service is currently down.=
</strong></p>
                  <p style=3D"line-height: 20px;">We will =
alert you when it&apos;s up again.</p>
                          </div>

            <div style=3D"padding: 20px; background-color: #f9f9f9; =
border-radius: 6px; margin-top: 25px; line-height: 16px;">

              <div style=3D"font-size: 12px; color: #687790;">Monitor =
name</div>
              <h2 style=3D"font-size: 14px; margin-bottom: 5px; =
margin-top: 3px;">{{.PipelineName}}</h2>
              <hr style=3D"border: 1px solid =
#dedede; border-bottom: 0; margin: 10px 0;">

              <div =
style=3D"color: #687790; font-size: 12px;">Checked URL</div>
              <h2 style=3D"font-size: 14px; margin-bottom: 5px; margin-top:=
 3px; line-height: 16px;"><code>
                <a href=3D"{{.Address}}" style=3D"color: #131a26 !important; text-decoration: none !=
important;">{{.Address}}</a></code></h2>
              <hr style=3D"border: 1px solid #dedede; border-bottom: 0; =
margin: 10px 0;">

<div style=3D"color: #687790; font-size: =
12px;">Datacenters Showed Problem</div>
<h2 style=3D"font-size: 14px; =
margin-bottom: 5px; margin-top: 3px;">
                  {{.Datacenters}}
              </h2>
              <hr style=3D"border: 1px solid #dedede; =
border-bottom: 0; margin: 10px 0;">


              <div style=3D"color: #687790; font-size: =
12px;">Root cause</div>
              <h2 style=3D"font-size: 14px; =
margin-bottom: 5px; margin-top: 3px;">
                  {{.RootCause}}
              </h2>
              <hr style=3D"border: 1px solid #dedede; =
border-bottom: 0; margin: 10px 0;">

              <div style=3D"font-size:=
 12px; color: #687790;">Incident resolved at</div>
              <h2 =
style=3D"font-size: 14px; margin-bottom: 5px; margin-top: 3px;">{{.Time}}</h2>
              <hr style=3D"border: 1px solid #dedede; =
border-bottom: 0; margin: 10px 0;">
            </div>
          </div>
        </div>
      </div>
  </body>
</html>`
