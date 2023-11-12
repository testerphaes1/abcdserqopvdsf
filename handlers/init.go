package handlers

import (
	"crypto/rsa"
)

var PrivateKey *rsa.PrivateKey

func init() {
	//_, vConfig, err := config.ViperConfig()
	//if err != nil {
	//	panic(err)
	//}
	//psqlDb, err := utils.PostgresConnection(vConfig.Database.Pslq.Host, vConfig.Database.Pslq.Port, vConfig.Database.Pslq.User, vConfig.Database.Pslq.Password, vConfig.Database.Pslq.Database, vConfig.Database.Pslq.Ssl)
	//if err != nil {
	//	panic(err)
	//}
	//
	//authInfoRepo := repos.NewAuthInfoRepositoryRepository(psqlDb)
	//privateStr, err := authInfoRepo.GetAuthKeys(context.TODO())
	//if err != nil {
	//	panic(err)
	//}
	//
	//block, _ := pem.Decode([]byte(privateStr))
	//if block == nil {
	//	panic("failed to parse PEM block containing the key")
	//}
	//PrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	//if err != nil {
	//	panic(err)
	//}
	//
	//psqlDb.Close()
}
