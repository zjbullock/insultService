package helper

import (
	"encoding/json"
	"github.com/juju/loggo"
	"github.com/pkg/errors"
	"insultService/model"
	"io/ioutil"
	"os"
)

func GetSecrets(log loggo.Logger) (*model.Secrets, error) {
	file, err := os.Open("./secrets/secrets.json")
	if err != nil {
		return nil, errors.Wrapf(err, "error getting twilio stuff")
	}
	defer file.Close()
	log.Infof("Helper - GetSecrets - Reading File")
	byteValue, _ := ioutil.ReadAll(file)

	log.Infof("Helper - GetSecrets - Now unmarshalling secrets")

	var secret model.Secrets
	err = json.Unmarshal(byteValue, &secret)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling json")
	}
	log.Infof("Helper - GetSecrets - Now returning with secrets")

	return &secret, nil
}
