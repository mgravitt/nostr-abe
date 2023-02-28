package main

import (
	"fmt"

	"github.com/fentec-project/gofe/abe"
	faker "github.com/go-faker/faker/v4"
)

type User struct {
	FirstName string `faker:"first_name"`
	Region    string `faker:"oneof: americas, europe, asia, australia, icelandic"`
	Job       string `faker:"oneof: plumber, electrician, scientist, programmer, nurse"`
	Tier      string `faker:"oneof: platinum, gold, silver, bronze"`
	RegionKey abe.MAABEKey
	JobKey    abe.MAABEKey
	TierKey   abe.MAABEKey
}

func NewUser(regionAuthority abe.MAABEAuth, jobAuthority abe.MAABEAuth, tierAuthority abe.MAABEAuth) User {

	user := User{}
	err := faker.FakeData(&user)
	if err != nil {
		panic(err)
	}

	user.Region = "region:" + user.Region
	user.Job = "job:" + user.Job
	user.Tier = "tier:" + user.Tier
	fmt.Println("\nUser		: ", user.FirstName)
	fmt.Println("-- Region	: ", user.Region)
	fmt.Println("-- Job		: ", user.Job)
	fmt.Println("-- Tier		: ", user.Tier)

	// authority 1 issues keys to user
	rKeys, err := regionAuthority.GenerateAttribKeys(user.FirstName, []string{user.Region})
	if err != nil {
		panic(err)
	}
	user.RegionKey = *rKeys[0]

	jobKeys, err := jobAuthority.GenerateAttribKeys(user.FirstName, []string{user.Job})
	if err != nil {
		panic(err)
	}
	user.JobKey = *jobKeys[0]

	tierKeys, err := tierAuthority.GenerateAttribKeys(user.FirstName, []string{user.Tier})
	if err != nil {
		panic(err)
	}
	user.TierKey = *tierKeys[0]
	return user
}

func (u *User) CanDecrypt(plaintext string, maabe *abe.MAABE, ct *abe.MAABECipher) bool {
	ks1 := []*abe.MAABEKey{&u.RegionKey, &u.JobKey, &u.TierKey}

	msg1, err := maabe.Decrypt(ct, ks1)
	if err != nil {
		return false
	}
	return plaintext == msg1
}

type Authority struct {
	label        string
	designations []string
	auth         abe.MAABEAuth
}

func main() {
	// create new MAABE struct with Global Parameters
	maabe := abe.NewMAABE()

	// create three authorities that independently have the ability to designate users as one of their values
	possibleRegions := []string{"region:americas", "region:europe", "region:asia", "region:australia", "region:icelandic"}
	possibleJobs := []string{"job:plumber", "job:electrician", "job:scientist", "job:programmer", "job:nurse"}
	possibleTiers := []string{"tier:platinum", "tier:gold", "tier:silver", "tier:bronze"}

	fmt.Println("Possible Regions	: ", possibleRegions)
	fmt.Println("Possible Jobs		: ", possibleJobs)
	fmt.Println("Possible Tiers		: ", possibleTiers)

	regionAuthority, err := maabe.NewMAABEAuth("region", possibleRegions)
	if err != nil {
		panic(err)
	}
	jobAuthority, err := maabe.NewMAABEAuth("job", possibleJobs)
	if err != nil {
		panic(err)
	}
	tierAuthority, err := maabe.NewMAABEAuth("tier", possibleTiers)
	if err != nil {
		panic(err)
	}

	// create a msp struct out of the boolean formula
	decryptionPolicy := "((region:americas AND job:scientist) OR (region:icelandic AND job:plumber)) OR (tier:platinum)"
	fmt.Println("\nDecryption Policy	: ", decryptionPolicy)
	msp, err := abe.BooleanToMSP(decryptionPolicy, false)
	if err != nil {
		panic(err)
	}

	// define the set of all public keys we use
	pks := []*abe.MAABEPubKey{regionAuthority.PubKeys(), jobAuthority.PubKeys(), tierAuthority.PubKeys()}

	// choose a message
	msg := "Attack at dawn!"

	// encrypt the message with the decryption policy in msp
	ct, err := maabe.Encrypt(msg, msp, pks)
	if err != nil {
		panic(err)
	}

	// generate and print 10 random users to see who can decrypt the cipher text based on
	// their designated attributes and the ciphertext
	for i := 0; i < 10; i++ {
		u := NewUser(*regionAuthority, *jobAuthority, *tierAuthority)
		fmt.Println("Can decrypt?	: ", u.CanDecrypt(msg, maabe, ct))
	}
}
