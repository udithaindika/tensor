package credentials

import (
	"bitbucket.pearson.com/apseng/tensor/models"
	"log"
	"gopkg.in/mgo.v2/bson"
	"bitbucket.pearson.com/apseng/tensor/db"
	"time"
)

// hideEncrypted is replace encrypted fields by $encrypted$
func hideEncrypted(c *models.Credential) {
	encrypted := "$encrypted$"
	c.Password = &encrypted
	c.SshKeyData = &encrypted
	c.SshKeyUnlock = &encrypted
	c.BecomePassword = &encrypted
	c.VaultPassword = &encrypted
	c.AuthorizePassword = &encrypted
}

func addActivity(crdID bson.ObjectId, userID bson.ObjectId, desc string) {

	c := db.C(db.ACTIVITY_STREAM)

	err := c.Insert(models.Activity{
		ID: bson.NewObjectId(),
		ActorID: userID,
		Type: _CTX_CREDENTIAL,
		ObjectID: crdID,
		Description: desc,
		Created: time.Now(),
	});

	if err != nil {
		log.Println("Failed to add new Activity", err)
	}
}