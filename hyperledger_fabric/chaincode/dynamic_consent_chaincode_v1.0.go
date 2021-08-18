© 2020-2021 vtw <hyunwoo.kwon@vtw.co.kr>
Copyright (c) 2020, vtw. co ltd 

package chaincode

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type UserInfo struct {
	Key          string `json:"key"`
	UserName     string `json:"userName"`
	UserBirthday string `json:"userBirthday"`
	UserSex      string `json:"userSex"`
}

type User struct {
	UserSeq  string `json:"userSeq"`
	Address  string `json:"address"`
	UserInfo string `json:"userInfo"`
	UseYn    bool   `json:"useYn"`
}

type UserRes struct {
	UserSeq string `json:"userSeq"`
}

type Org struct {
	OrgSeq     string `json:"orgSeq"`
	OrgAddress string `json:"orgAddress"`
	Key        string `json:"key"`
	OrgName    string `json:"orgName"`
	OrgType    string `json:"orgType"`
	UseYn      bool   `json:"useYn"`
}

type OrgRes struct {
	OrgSeq string `json:"orgSeq"`
}

type OrgUser struct {
	OrgUserSeq   string `json:"orgUserSeq"`
	OrgAddress   string `json:"orgAddress"`
	UserAddress  string `json:"userAddress"`
	UserInfoHash string `json:"userInfoHash"`
	OrgID        string `json:"orgID"`
}

type OrgUserHash struct {
	OrgID        string `json:"orgID"`
	UserName     string `json:"userName"`
	UserBirthday string `json:"userBirthday"`
	UserSex      string `json:"userSex"`
}

type OrgUserRes struct {
	OrgUserSeq string `json:"orgUserSeq"`
}

type Ad struct {
	Address  string `json:"address"`
	Document string `json:"document"`
	UseYn    bool   `json:"useYn"`
}

type AdRes struct {
	AdSeq string `json:"adSeq"`
}

type AdReadRes struct {
	Result bool `json:"result"`
	UseYn  bool `json:"useYn"`
}

type Agree struct {
	AgreeSeq string `json:"agreeSeq"`
	Address  string `json:"address"`
	AdSeq    string `json:"adSeq"`
	AgreeYn  bool   `json:"agreeYn"`
}

type AgreeRes struct {
	AgreeSeq string `json:"agreeSeq"`
}

type AgreeReadRes struct {
	Result  bool `json:"result"`
	AgreeYn bool `json:"agreeYn"`
}

type CategoryKey struct {
	Key string `json:"key"`
	Idx int    `json:"idx"`
}

type Medi struct {
	MediSeq  string `json:"mediSeq"`
	UserSeq  string `json:"userSeq"`
	AgreeSeq string `json:"agreeSeq"`
	Status   string `json:"status"`
}

type MediRes struct {
	MediSeq string `json:"mediSeq"`
}

type MediCompareRes struct {
	Result bool `json:"result"`
}

type MediDetail struct {
	MediDetailSeq string `json:"mediDetailSeq"`
	Address       string `json:"address"`
	AgreeSeq      string `json:"agreeSeq"`
	MediSeq       string `json:"mediSeq"`
	MediInfo      string `json:"mediInfo"`
	Status        string `json:"status"`
}

type MediDetailRes struct {
	MediDetailSeq string `json:"mediDetailSeq"`
}

func generateKey(ctx contractapi.TransactionContextInterface, key string, category string) ([]byte, error) {

	var isFirst bool = false

	keyJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return keyJSON, errors.New("Error to get key")
	}

	categoryKey := CategoryKey{}
	json.Unmarshal(keyJSON, &categoryKey)
	var tempIdx string
	tempIdx = strconv.Itoa(categoryKey.Idx)
	fmt.Println(categoryKey)
	fmt.Println("Key is " + strconv.Itoa(len(categoryKey.Key)))
	if len(categoryKey.Key) == 0 || categoryKey.Key == "" {
		isFirst = true
		categoryKey.Key = category
	}
	if !isFirst {
		categoryKey.Idx = categoryKey.Idx + 1
	}

	fmt.Println("Last CategoryKey is " + categoryKey.Key + " : " + tempIdx)

	categoryBytes, _ := json.Marshal(categoryKey)

	return categoryBytes, err
}

func getAddress(privateKeyEncoded string) (string, error) {

	if privateKeyEncoded[:2] == "0x" {
		privateKeyEncoded = privateKeyEncoded[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyEncoded)
	if err != nil {
		return "", err
	}

	// create public key
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Printf("Public key:\t %s\n", hexutil.Encode(publicKeyBytes)[4:])

	// create address
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()[2:]
	fmt.Printf("Public address (from ECDSA): \t%s\n", address)

	return address, err
}

func findUser(ctx contractapi.TransactionContextInterface, address string, hash string) (*User, error) {

	// Find latestKeyUser
	userKeyJSON, _ := ctx.GetStub().GetState("latestKeyUser")
	userKey := CategoryKey{}
	json.Unmarshal(userKeyJSON, &userKey)

	idxStr := strconv.Itoa(userKey.Idx + 1)

	var startKey = "USER0"
	var endKey = userKey.Key + idxStr

	resultUsersIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, errors.New("Error to get USER key range")
	}
	defer resultUsersIter.Close()

	for resultUsersIter.HasNext() {
		resultUser, _ := resultUsersIter.Next()

		userEach := User{}
		json.Unmarshal(resultUser.Value, &userEach)
		if hash == "" {
			if userEach.Address == address {
				return &userEach, err
			}
		} else if address == "" {
			if userEach.UserInfo == hash {
				return &userEach, err
			}
		} else if hash != "" && address != "" {
			if userEach.Address == address && userEach.UserInfo == hash {
				return &userEach, err
			}
		}
	}

	return nil, errors.New("Fail to find user address")
}

func checkUser(ctx contractapi.TransactionContextInterface, userAddress string, userInfo string) error {

	// Find latestKeyUser
	userKeyJSON, _ := ctx.GetStub().GetState("latestKeyUser")
	userKey := CategoryKey{}
	json.Unmarshal(userKeyJSON, &userKey)

	idxStr := strconv.Itoa(userKey.Idx + 1)

	var startKey = "ORG0"
	var endKey = userKey.Key + idxStr

	resultUsersIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return errors.New("Error to get ORG key range")
	}
	defer resultUsersIter.Close()

	for resultUsersIter.HasNext() {
		resultUser, _ := resultUsersIter.Next()

		userEach := User{}
		json.Unmarshal(resultUser.Value, &userEach)

		if userEach.Address == userAddress {
			return errors.New("address already exists.")
		} else if userEach.UserInfo == userInfo {
			return errors.New("userInfo already exists.")
		}
	}

	return nil
}

func findReissueUser(ctx contractapi.TransactionContextInterface, hash string) (*User, error) {

	// Find latestKeyUser
	userKeyJSON, _ := ctx.GetStub().GetState("latestKeyUser")
	userKey := CategoryKey{}
	json.Unmarshal(userKeyJSON, &userKey)

	idxStr := strconv.Itoa(userKey.Idx + 1)

	var startKey = "USER0"
	var endKey = userKey.Key + idxStr

	resultUsersIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, errors.New("Error to get USER key range")
	}
	defer resultUsersIter.Close()

	for resultUsersIter.HasNext() {
		resultUser, _ := resultUsersIter.Next()

		userEach := User{}
		json.Unmarshal(resultUser.Value, &userEach)

		if userEach.UserInfo == hash && !userEach.UseYn {
			return &userEach, err
		} else if userEach.UserInfo == hash && userEach.UseYn {
			return nil, errors.New("UseYn is true")
		}

	}

	return nil, errors.New("Fail to find user address")
}

func findOrg(ctx contractapi.TransactionContextInterface, address string) (*Org, error) {

	// Find latestKeyUser
	orgKeyJSON, _ := ctx.GetStub().GetState("latestKeyOrg")
	orgKey := CategoryKey{}
	json.Unmarshal(orgKeyJSON, &orgKey)

	idxStr := strconv.Itoa(orgKey.Idx + 1)

	var startKey = "ORG0"
	var endKey = orgKey.Key + idxStr

	resultOrgsIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, errors.New("Error to get ORG key range")
	}
	defer resultOrgsIter.Close()

	for resultOrgsIter.HasNext() {
		resultOrg, _ := resultOrgsIter.Next()

		orgEach := Org{}
		json.Unmarshal(resultOrg.Value, &orgEach)

		if orgEach.OrgAddress == address {
			return &orgEach, err
		}

	}
	return nil, errors.New("Fail to find org address")
}

func findReissueOrg(ctx contractapi.TransactionContextInterface, key string) (*Org, error) {

	// Find latestKeyUser
	orgKeyJSON, _ := ctx.GetStub().GetState("latestKeyOrg")
	orgKey := CategoryKey{}
	json.Unmarshal(orgKeyJSON, &orgKey)

	idxStr := strconv.Itoa(orgKey.Idx + 1)

	var startKey = "ORG0"
	var endKey = orgKey.Key + idxStr

	resultOrgsIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, errors.New("Error to get ORG key range")
	}
	defer resultOrgsIter.Close()

	for resultOrgsIter.HasNext() {
		resultOrg, _ := resultOrgsIter.Next()

		orgEach := Org{}
		json.Unmarshal(resultOrg.Value, &orgEach)

		if orgEach.Key == key && !orgEach.UseYn {
			return &orgEach, err
		} else if orgEach.Key == key && orgEach.UseYn {
			return nil, errors.New("UseYn is true")
		}

	}
	return nil, errors.New("Fail to find org address")
}

func checkOrg(ctx contractapi.TransactionContextInterface, key string, orgType string) error {

	// Find latestKeyUser
	orgKeyJSON, _ := ctx.GetStub().GetState("latestKeyOrg")
	orgKey := CategoryKey{}
	json.Unmarshal(orgKeyJSON, &orgKey)

	idxStr := strconv.Itoa(orgKey.Idx + 1)

	var startKey = "ORG0"
	var endKey = orgKey.Key + idxStr

	resultOrgsIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return errors.New("Error to get ORG key range")
	}
	defer resultOrgsIter.Close()

	for resultOrgsIter.HasNext() {
		resultOrg, _ := resultOrgsIter.Next()

		orgEach := Org{}
		json.Unmarshal(resultOrg.Value, &orgEach)

		if orgEach.OrgType == "super" && orgType == "super" {
			return errors.New("Super admin already exists.")
		} else if orgEach.Key == key {
			return errors.New("address already exists.")
		}
	}

	return nil
}

func findOrgUser(ctx contractapi.TransactionContextInterface, orgAddress string, orgUserHash string) (*OrgUser, error) {

	// Find latestKeyOrgUser
	orgUserKeyJSON, _ := ctx.GetStub().GetState("latestKeyOrgUser")
	orgUserKey := CategoryKey{}
	json.Unmarshal(orgUserKeyJSON, &orgUserKey)

	idxStr := strconv.Itoa(orgUserKey.Idx + 1)

	var startKey = "ORGUSER0"
	var endKey = orgUserKey.Key + idxStr

	resultOrgUsersIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, errors.New("Error to get ORGUSER key range")
	}
	defer resultOrgUsersIter.Close()

	for resultOrgUsersIter.HasNext() {
		resultOrgUser, _ := resultOrgUsersIter.Next()

		orgUserEach := OrgUser{}
		json.Unmarshal(resultOrgUser.Value, &orgUserEach)

		if orgUserEach.OrgAddress == orgAddress && orgUserEach.UserInfoHash == orgUserHash {
			return &orgUserEach, err
		}
	}

	return nil, errors.New("Fail to find ORGUSER")
}

func changeAddressOrgUserOrg(ctx contractapi.TransactionContextInterface, beforeAddress string, afterAddress string) error {

	// Find latestKeyUser
	userKeyJSON, _ := ctx.GetStub().GetState("latestKeyOrgUser")
	userKey := CategoryKey{}
	json.Unmarshal(userKeyJSON, &userKey)

	idxStr := strconv.Itoa(userKey.Idx + 1)

	var startKey = "ORGUSER0"
	var endKey = userKey.Key + idxStr

	resultIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return errors.New("Error to get key range")
	}
	defer resultIter.Close()

	for resultIter.HasNext() {
		result, _ := resultIter.Next()

		orgUserEach := OrgUser{}
		json.Unmarshal(result.Value, &orgUserEach)

		if orgUserEach.OrgAddress == beforeAddress {
			orgUserEach.OrgAddress = afterAddress

			orgUserEachJSON, err := json.Marshal(orgUserEach)
			if err != nil {
				return errors.New("Error to convert marshal")
			}

			ctx.GetStub().PutState(result.Key, orgUserEachJSON)
		}
	}
	return nil
}

func changeAddressOrgUserUser(ctx contractapi.TransactionContextInterface, beforeAddress string, afterAddress string) error {

	// Find latestKeyUser
	userKeyJSON, _ := ctx.GetStub().GetState("latestKeyOrgUser")
	userKey := CategoryKey{}
	json.Unmarshal(userKeyJSON, &userKey)

	idxStr := strconv.Itoa(userKey.Idx + 1)

	var startKey = "ORGUSER0"
	var endKey = userKey.Key + idxStr

	resultIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return errors.New("Error to get key range")
	}
	defer resultIter.Close()

	for resultIter.HasNext() {
		result, _ := resultIter.Next()

		orgUserEach := OrgUser{}
		json.Unmarshal(result.Value, &orgUserEach)

		if orgUserEach.UserAddress == beforeAddress {
			orgUserEach.UserAddress = afterAddress

			orgUserEachJSON, err := json.Marshal(orgUserEach)
			if err != nil {
				return errors.New("Error to convert marshal")
			}

			ctx.GetStub().PutState(result.Key, orgUserEachJSON)
		}
	}
	return nil
}

func changeAddressAd(ctx contractapi.TransactionContextInterface, beforeAddress string, afterAddress string) error {

	// Find latestKeyUser
	userKeyJSON, _ := ctx.GetStub().GetState("latestKeyAd")
	userKey := CategoryKey{}
	json.Unmarshal(userKeyJSON, &userKey)

	idxStr := strconv.Itoa(userKey.Idx + 1)

	var startKey = "AD0"
	var endKey = userKey.Key + idxStr

	resultIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return errors.New("Error to get key range")
	}
	defer resultIter.Close()

	for resultIter.HasNext() {
		result, _ := resultIter.Next()

		adEach := Ad{}
		json.Unmarshal(result.Value, &adEach)

		if adEach.Address == beforeAddress {
			adEach.Address = afterAddress

			adEachJSON, err := json.Marshal(adEach)
			if err != nil {
				return errors.New("Error to convert marshal")
			}

			ctx.GetStub().PutState(result.Key, adEachJSON)
		}
	}
	return nil
}

func changeAddressAgree(ctx contractapi.TransactionContextInterface, beforeAddress string, afterAddress string) error {

	// Find latestKeyUser
	userKeyJSON, _ := ctx.GetStub().GetState("latestKeyAgree")
	userKey := CategoryKey{}
	json.Unmarshal(userKeyJSON, &userKey)

	idxStr := strconv.Itoa(userKey.Idx + 1)

	var startKey = "AGREE0"
	var endKey = userKey.Key + idxStr

	resultIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return errors.New("Error to get key range")
	}
	defer resultIter.Close()

	for resultIter.HasNext() {
		result, _ := resultIter.Next()

		agreeEach := Agree{}
		json.Unmarshal(result.Value, &agreeEach)

		if agreeEach.Address == beforeAddress {
			agreeEach.Address = afterAddress

			agreeEachJSON, err := json.Marshal(agreeEach)
			if err != nil {
				return errors.New("Error to convert marshal")
			}

			ctx.GetStub().PutState(result.Key, agreeEachJSON)
		}
	}
	return nil
}

func changeAddressMediDetail(ctx contractapi.TransactionContextInterface, beforeAddress string, afterAddress string) error {

	// Find latestKeyUser
	userKeyJSON, _ := ctx.GetStub().GetState("latestKeyMediDetail")
	userKey := CategoryKey{}
	json.Unmarshal(userKeyJSON, &userKey)

	idxStr := strconv.Itoa(userKey.Idx + 1)

	var startKey = "MEDIDETAIL0"
	var endKey = userKey.Key + idxStr

	resultIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return errors.New("Error to get key range")
	}
	defer resultIter.Close()

	for resultIter.HasNext() {
		result, _ := resultIter.Next()

		mediDetailEach := MediDetail{}
		json.Unmarshal(result.Value, &mediDetailEach)

		if mediDetailEach.Address == beforeAddress {
			mediDetailEach.Address = afterAddress

			mediDetailEachJSON, err := json.Marshal(mediDetailEach)
			if err != nil {
				return errors.New("Error to convert marshal")
			}

			ctx.GetStub().PutState(result.Key, mediDetailEachJSON)
		}
	}
	return nil
}

func findAd(ctx contractapi.TransactionContextInterface, address string, docMdStr string) (*Ad, string, error) {

	// Find latestKeyAd
	adKeyJSON, _ := ctx.GetStub().GetState("latestKeyAd")
	adKey := CategoryKey{}
	json.Unmarshal(adKeyJSON, &adKey)

	idxStr := strconv.Itoa(adKey.Idx + 1)

	var startKey = "AD0"
	var endKey = adKey.Key + idxStr

	resultAdsIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, "", errors.New("Error to get AD key range")
	}
	defer resultAdsIter.Close()

	for resultAdsIter.HasNext() {
		resultAd, _ := resultAdsIter.Next()

		adEach := Ad{}
		json.Unmarshal(resultAd.Value, &adEach)

		if adEach.Address == address && adEach.Document == docMdStr {
			return &adEach, resultAd.Key, err
		}
	}

	return nil, "", errors.New("Fail to find AD")
}

func findAgree(ctx contractapi.TransactionContextInterface, address string, adSeq string) (*Agree, string, error) {

	// Find latestKeyAgree
	agreeKeyJSON, _ := ctx.GetStub().GetState("latestKeyAgree")
	agreeKey := CategoryKey{}
	json.Unmarshal(agreeKeyJSON, &agreeKey)

	idxStr := strconv.Itoa(agreeKey.Idx + 1)

	var startKey = "AGREE0"
	var endKey = agreeKey.Key + idxStr

	resultAgreesIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, "", errors.New("Error to get AGREE key range")
	}
	defer resultAgreesIter.Close()

	for resultAgreesIter.HasNext() {
		resultAgree, _ := resultAgreesIter.Next()

		agreeEach := Agree{}
		json.Unmarshal(resultAgree.Value, &agreeEach)
		if agreeEach.Address == address && agreeEach.AdSeq == adSeq {
			return &agreeEach, resultAgree.Key, err
		}
	}

	return nil, "", errors.New("Fail to find AGREE")
}

func findMedi(ctx contractapi.TransactionContextInterface, orgAddress string, agreeSeq string) ([]Medi, error) {

	mediSeqList := []Medi{}

	// Find latestKeyMedi
	mediKeyJSON, _ := ctx.GetStub().GetState("latestKeyMedi")
	mediKey := CategoryKey{}
	json.Unmarshal(mediKeyJSON, &mediKey)

	idxStr := strconv.Itoa(mediKey.Idx + 1)

	var startKey = "MEDI0"
	var endKey = mediKey.Key + idxStr

	resultMedisIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, errors.New("Error to get MEDI key range")
	}
	defer resultMedisIter.Close()

	for resultMedisIter.HasNext() {
		resultMedi, _ := resultMedisIter.Next()

		mediEach := Medi{}
		json.Unmarshal(resultMedi.Value, &mediEach)

		var agreeJSON []byte

		if agreeSeq == "" {
			// MEDI 안에 AGREE 검색 후 ORG address 추출
			agreeJSON, err = ctx.GetStub().GetState(mediEach.AgreeSeq)
		} else if agreeSeq == mediEach.AgreeSeq {
			// 동의정보 seq 쿼리
			agreeJSON, err = ctx.GetStub().GetState(agreeSeq)
		} else {
			continue
		}

		if err != nil {
			return nil, err
		}
		agree := Agree{}
		json.Unmarshal(agreeJSON, &agree)

		// 조직 seq 쿼리
		adJSON, err := ctx.GetStub().GetState(agree.AdSeq)
		if err != nil {
			return nil, err
		}

		ad := Ad{}
		json.Unmarshal(adJSON, &ad)

		if ad.Address == orgAddress && mediEach.Status == "request" {
			mediSeqList = append(mediSeqList, mediEach)
		}

	}

	return mediSeqList, err
}

func findMediSearch(ctx contractapi.TransactionContextInterface, userSeq string, orgSeq string, status string) ([]Medi, error) {

	mediSeqList := []Medi{}

	// Find latestKeyMedi
	mediKeyJSON, _ := ctx.GetStub().GetState("latestKeyMedi")
	mediKey := CategoryKey{}
	json.Unmarshal(mediKeyJSON, &mediKey)

	idxStr := strconv.Itoa(mediKey.Idx + 1)

	var startKey = "MEDI0"
	var endKey = mediKey.Key + idxStr

	resultMedisIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, errors.New("Error to get MEDI key range")
	}
	defer resultMedisIter.Close()

	for resultMedisIter.HasNext() {
		resultMedi, _ := resultMedisIter.Next()

		mediEach := Medi{}
		json.Unmarshal(resultMedi.Value, &mediEach)

		// 동의정보 seq 쿼리
		agreeJSON, err := ctx.GetStub().GetState(mediEach.AgreeSeq)
		if err != nil {
			return nil, errors.New("agreement sequence not found")
		}
		agree := Agree{}
		json.Unmarshal(agreeJSON, &agree)

		// 동의서 seq 쿼리
		adJSON, err := ctx.GetStub().GetState(agree.AdSeq)
		if err != nil {
			return nil, err
		}

		ad := Ad{}
		json.Unmarshal(adJSON, &ad)

		if userSeq == "" {
			if orgSeq == "" {
				if status == "" {
					return nil, errors.New("Enter at least one of the three parameters without privateKey")
				} else {
					if mediEach.Status == status {
						mediSeqList = append(mediSeqList, mediEach)
					}
				}
			} else {

				// 기관 찾기
				orgJSON, err := ctx.GetStub().GetState(orgSeq)
				if err != nil {
					return nil, errors.New("org not found")
				}

				org := Org{}
				json.Unmarshal(orgJSON, &org)

				if status == "" {
					if ad.Address == org.OrgAddress {
						mediSeqList = append(mediSeqList, mediEach)
					}
				} else {
					if ad.Address == org.OrgAddress && mediEach.Status == status {
						mediSeqList = append(mediSeqList, mediEach)
					}
				}
			}
		} else {
			if orgSeq == "" {
				if status == "" {
					if mediEach.UserSeq == userSeq {
						mediSeqList = append(mediSeqList, mediEach)
					}
				} else {
					if mediEach.UserSeq == userSeq && mediEach.Status == status {
						mediSeqList = append(mediSeqList, mediEach)
					}
				}
			} else {

				// 기관 찾기
				orgJSON, err := ctx.GetStub().GetState(orgSeq)
				if err != nil {
					return nil, errors.New("org not found")
				}

				org := Org{}
				json.Unmarshal(orgJSON, &org)

				if status == "" {
					if mediEach.UserSeq == userSeq && ad.Address == org.OrgAddress {
						mediSeqList = append(mediSeqList, mediEach)
					}
				} else {
					if mediEach.UserSeq == userSeq && ad.Address == org.OrgAddress && mediEach.Status == status {
						mediSeqList = append(mediSeqList, mediEach)
					}
				}
			}
		}

	}

	return mediSeqList, err
}

func findMediDetail(ctx contractapi.TransactionContextInterface, mediSeq string, agreeSeq string) (*MediDetail, error) {

	// Find latestKeyMediDetail
	mediDetailKeyJSON, _ := ctx.GetStub().GetState("latestKeyMediDetail")
	mediDetailKey := CategoryKey{}
	json.Unmarshal(mediDetailKeyJSON, &mediDetailKey)

	idxStr := strconv.Itoa(mediDetailKey.Idx + 1)

	var startKey = "MEDIDETAIL0"
	var endKey = mediDetailKey.Key + idxStr

	resultMediDetailsIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, errors.New("Error to get MEDIDETAIL key range")
	}
	defer resultMediDetailsIter.Close()

	for resultMediDetailsIter.HasNext() {
		resultMediDetail, _ := resultMediDetailsIter.Next()

		mediDetailEach := MediDetail{}
		json.Unmarshal(resultMediDetail.Value, &mediDetailEach)

		if mediDetailEach.MediSeq == mediSeq && mediDetailEach.AgreeSeq == agreeSeq {
			return &mediDetailEach, err
		}
	}

	return nil, errors.New("Fail to find MEDIDETAIL")
}

func findMediDetailSearch(ctx contractapi.TransactionContextInterface, address string, mediDetailSeq string) ([]MediDetail, error) {

	mediDetailList := []MediDetail{}

	// 기관 찾기
	mediDetailJSON, err := ctx.GetStub().GetState(mediDetailSeq)
	if err != nil {
		return nil, err
	}

	mediDetail := MediDetail{}
	json.Unmarshal(mediDetailJSON, &mediDetail)

	if address != "" {
		if mediDetail.Address == address {
			mediDetailList = append(mediDetailList, mediDetail)
			return mediDetailList, err
		} else {
			return nil, errors.New("Mismatch Org Address")
		}
	}

	// Find latestKeyMedi
	mediKeyJSON, _ := ctx.GetStub().GetState("latestKeyMediDetail")
	mediKey := CategoryKey{}
	json.Unmarshal(mediKeyJSON, &mediKey)

	idxStr := strconv.Itoa(mediKey.Idx + 1)

	var startKey = "MEDIDETAIL0"
	var endKey = mediKey.Key + idxStr

	resultMediDetailsIter, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, errors.New("Error to get MEDI key range")
	}
	defer resultMediDetailsIter.Close()

	for resultMediDetailsIter.HasNext() {
		resultMediDetail, _ := resultMediDetailsIter.Next()

		mediDetailEach := MediDetail{}
		json.Unmarshal(resultMediDetail.Value, &mediDetailEach)

		if mediDetailEach.Address == mediDetail.Address {
			mediDetailList = append(mediDetailList, mediDetailEach)
		}
	}

	return mediDetailList, err
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// 유저 생성
func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, key string, userName string, userBirthday string, userSex string) (*UserRes, error) {

	var categoryKey = CategoryKey{}
	categoryBytes, err := generateKey(ctx, "latestKeyUser", "USER")
	if err != nil {
		return nil, err
	}
	json.Unmarshal(categoryBytes, &categoryKey)
	keyidx := strconv.Itoa(categoryKey.Idx)
	var keyString = categoryKey.Key + keyidx

	userAddress, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}

	userInfo := UserInfo{
		Key:          key,
		UserName:     userName,
		UserBirthday: userBirthday,
		UserSex:      userSex,
	}

	// 유저 정보 해시화
	userInfoBytes, err := json.Marshal(userInfo)

	userInfohash := sha256.New()
	userInfohash.Write(userInfoBytes)

	md := userInfohash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	// 유저 찾기
	err = checkUser(ctx, userAddress, mdStr)
	if err != nil {
		return nil, err
	}

	user := User{
		UserSeq:  keyString,
		Address:  userAddress,
		UserInfo: mdStr,
		UseYn:    true,
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(keyString, userJSON)
	if err != nil {
		return nil, err
	}

	categoryKeyJSON, err := json.Marshal(categoryKey)
	if err != nil {
		return nil, err
	}
	err = ctx.GetStub().PutState("latestKeyUser", categoryKeyJSON)
	if err != nil {
		return nil, err
	}

	userRes := UserRes{
		UserSeq: keyString,
	}

	return &userRes, nil
}

// 유저 재발급
func (s *SmartContract) ReissueUser(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, key string, userName string, userBirthday string, userSex string) (*UserRes, error) {

	userAddress, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}
	// 유저 정보
	userInfo := UserInfo{
		Key:          key,
		UserName:     userName,
		UserBirthday: userBirthday,
		UserSex:      userSex,
	}

	// 유저 정보 해시화
	userInfoBytes, err := json.Marshal(userInfo)

	userInfohash := sha256.New()
	userInfohash.Write(userInfoBytes)

	md := userInfohash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	// 유저 찾기
	reissueUser, err := findReissueUser(ctx, mdStr)

	user := User{
		UserSeq:  reissueUser.UserSeq,
		Address:  userAddress,
		UserInfo: mdStr,
		UseYn:    true,
	}

	userJSON, err := json.Marshal(user)

	err = ctx.GetStub().PutState(reissueUser.UserSeq, userJSON)
	if err != nil {
		return nil, err
	}

	// OrgUser address 변경
	err = changeAddressOrgUserUser(ctx, reissueUser.Address, userAddress)
	if err != nil {
		return nil, err
	}

	// Agree address 변경
	err = changeAddressAgree(ctx, reissueUser.Address, userAddress)
	if err != nil {
		return nil, err
	}

	userRes := UserRes{
		UserSeq: reissueUser.UserSeq,
	}

	return &userRes, nil
}

// 기관 생성
func (s *SmartContract) CreateOrg(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, key string, orgName string, orgType string) (*OrgRes, error) {

	if orgType != "general" && orgType != "super" {
		return nil, errors.New("Error orgType")
	}

	var categoryKey = CategoryKey{}
	categoryBytes, err := generateKey(ctx, "latestKeyOrg", "ORG")
	if err != nil {
		return nil, err
	}
	json.Unmarshal(categoryBytes, &categoryKey)
	keyidx := strconv.Itoa(categoryKey.Idx)
	var keyString = categoryKey.Key + keyidx

	orgAddress, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}

	err = checkOrg(ctx, key, orgType)
	if err != nil {
		return nil, err
	}

	org := Org{
		OrgSeq:     keyString,
		OrgAddress: orgAddress,
		Key:        key,
		OrgName:    orgName,
		OrgType:    orgType,
		UseYn:      true,
	}

	orgJSON, err := json.Marshal(org)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(keyString, orgJSON)
	if err != nil {
		return nil, err
	}

	categoryKeyJSON, err := json.Marshal(categoryKey)
	if err != nil {
		return nil, err
	}

	// Org 최근키 등록
	err = ctx.GetStub().PutState("latestKeyOrg", categoryKeyJSON)
	if err != nil {
		return nil, err
	}

	orgRes := OrgRes{
		OrgSeq: keyString,
	}

	return &orgRes, nil
}

// 기관 재발급
func (s *SmartContract) ReissueOrg(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, key string, orgName string, orgType string) (*OrgRes, error) {

	if orgType != "general" {
		return nil, errors.New("Error orgType")
	}

	orgAddress, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}

	reissueOrg, err := findReissueOrg(ctx, key)
	if err != nil {
		return nil, err
	}

	org := Org{
		OrgSeq:     reissueOrg.OrgSeq,
		OrgAddress: orgAddress,
		Key:        key,
		OrgName:    orgName,
		OrgType:    orgType,
		UseYn:      true,
	}

	orgJSON, err := json.Marshal(org)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(reissueOrg.OrgSeq, orgJSON)
	if err != nil {
		return nil, err
	}

	// OrgUser address 변경
	err = changeAddressOrgUserOrg(ctx, reissueOrg.OrgAddress, orgAddress)
	if err != nil {
		return nil, err
	}

	// Ad address 변경
	err = changeAddressAd(ctx, reissueOrg.OrgAddress, orgAddress)
	if err != nil {
		return nil, err
	}

	// Medidetail address 변경
	err = changeAddressMediDetail(ctx, reissueOrg.OrgAddress, orgAddress)
	if err != nil {
		return nil, err
	}

	orgRes := OrgRes{
		OrgSeq: reissueOrg.OrgSeq,
	}

	return &orgRes, nil
}

// 유저 재발급
func (s *SmartContract) LockUser(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, userSeq string) (*User, error) {

	// 슈퍼 권한 체크
	orgAddress, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}
	orgInfo, err := findOrg(ctx, orgAddress)

	if err != nil {
		return nil, err
	}

	if orgInfo.OrgType != "super" {
		return nil, err
	}

	// 유저 정보 seq 쿼리
	userJSON, err := ctx.GetStub().GetState(userSeq)
	if err != nil {
		return nil, err
	}
	user := User{}
	json.Unmarshal(userJSON, &user)

	user.UseYn = false
	userJSON, err = json.Marshal(user)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(user.UserSeq, userJSON)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// 기관 재발급
func (s *SmartContract) LockOrg(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, orgSeq string) (*Org, error) {

	// 슈퍼 권한 체크
	orgAddress, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}
	orgInfo, err := findOrg(ctx, orgAddress)

	if err != nil {
		return nil, err
	}

	if orgInfo.OrgType != "super" {
		return nil, err
	}

	// 유저 정보 seq 쿼리
	orgJSON, err := ctx.GetStub().GetState(orgSeq)
	if err != nil {
		return nil, err
	}
	org := Org{}
	json.Unmarshal(orgJSON, &org)

	org.UseYn = false
	orgJSON, err = json.Marshal(org)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(org.OrgSeq, orgJSON)
	if err != nil {
		return nil, err
	}

	return &org, nil
}

// 기관에 유저 등록
func (s *SmartContract) CreateOrgUser(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, orgSeq string, orgID string, userName string, userBirthday string, userSex string) (*OrgUserRes, error) {

	var categoryKey = CategoryKey{}
	categoryBytes, err := generateKey(ctx, "latestKeyOrgUser", "ORGUSER")
	if err != nil {
		return nil, err
	}
	json.Unmarshal(categoryBytes, &categoryKey)
	keyidx := strconv.Itoa(categoryKey.Idx)
	var keyString = categoryKey.Key + keyidx

	// 조직 seq 쿼리
	orgJSON, err := ctx.GetStub().GetState(orgSeq)
	if err != nil {
		return nil, err
	}

	// 주소 찾기
	userAddress, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}

	// 유저 찾기
	user, err := findUser(ctx, userAddress, "")
	if err != nil {
		return nil, err
	}

	orgUserHash := OrgUserHash{
		OrgID:        orgID,
		UserName:     userName,
		UserBirthday: userBirthday,
		UserSex:      userSex,
	}

	// 유저 정보 해시화
	orgUserHashBytes, err := json.Marshal(orgUserHash)

	orgUserHashhash := sha256.New()
	orgUserHashhash.Write(orgUserHashBytes)

	md := orgUserHashhash.Sum(nil)
	orgUserHashMdStr := hex.EncodeToString(md)

	org := Org{}
	json.Unmarshal(orgJSON, &org)

	orgUser := OrgUser{
		OrgUserSeq:   keyString,
		OrgAddress:   org.OrgAddress,
		UserAddress:  user.Address,
		UserInfoHash: orgUserHashMdStr,
		OrgID:        orgID,
	}

	orgUserInfoJSON, err := json.Marshal(orgUser)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(keyString, orgUserInfoJSON)
	if err != nil {
		return nil, err
	}

	categoryKeyJSON, err := json.Marshal(categoryKey)
	if err != nil {
		return nil, err
	}
	err = ctx.GetStub().PutState("latestKeyOrgUser", categoryKeyJSON)
	if err != nil {
		return nil, err
	}

	orgUserRes := OrgUserRes{
		OrgUserSeq: keyString,
	}

	return &orgUserRes, nil
}

// 기관에 유저 등록
func (s *SmartContract) ReadUser(ctx contractapi.TransactionContextInterface, key string, userName string, userBirthday string, userSex string) (*User, error) {

	userInfo := UserInfo{
		Key:          key,
		UserName:     userName,
		UserBirthday: userBirthday,
		UserSex:      userSex,
	}

	// 유저 정보 해시화
	userInfoBytes, err := json.Marshal(userInfo)
	if err != nil {
		return nil, err
	}

	userInfohash := sha256.New()
	userInfohash.Write(userInfoBytes)

	md := userInfohash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	// 유저 찾기
	user, err := findUser(ctx, "", mdStr)
	if err != nil {
		return nil, err
	}

	// 유저 찾기
	return user, err

}

// 기관에 등록된 유저 검색
func (s *SmartContract) ReadOrgUser(ctx contractapi.TransactionContextInterface, orgSeq string, orgID string, userName string, userBirthday string, userSex string) (*OrgUser, error) {

	orgUserHash := OrgUserHash{
		OrgID:        orgID,
		UserName:     userName,
		UserBirthday: userBirthday,
		UserSex:      userSex,
	}

	// 조직내 가입한 유저 정보 해시화
	orgUserHashBytes, err := json.Marshal(orgUserHash)

	orgUserHashhash := sha256.New()
	orgUserHashhash.Write(orgUserHashBytes)

	md := orgUserHashhash.Sum(nil)
	orgUserHashMdStr := hex.EncodeToString(md)

	// 기관 찾기
	orgJSON, err := ctx.GetStub().GetState(orgSeq)
	if err != nil {
		return nil, err
	}

	org := Org{}
	json.Unmarshal(orgJSON, &org)

	// 유저 찾기
	return findOrgUser(ctx, org.OrgAddress, orgUserHashMdStr)

}

//// 동의서

// 동의서 등록
func (s *SmartContract) CreateDoc(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, document string, useYn bool) (*AdRes, error) {

	var categoryKey = CategoryKey{}
	categoryBytes, err := generateKey(ctx, "latestKeyAd", "AD")
	if err != nil {
		return nil, err
	}
	json.Unmarshal(categoryBytes, &categoryKey)
	keyidx := strconv.Itoa(categoryKey.Idx)
	var keyString = categoryKey.Key + keyidx

	address, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}

	org, err := findOrg(ctx, address)
	if err != nil {
		return nil, err
	}

	// 동의문서 해시화
	userInfohash := sha256.New()
	userInfohash.Write([]byte(document))

	md := userInfohash.Sum(nil)
	docMdStr := hex.EncodeToString(md)

	ad := Ad{
		Address:  org.OrgAddress,
		Document: docMdStr,
		UseYn:    useYn,
	}

	adJSON, err := json.Marshal(ad)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(keyString, adJSON)
	if err != nil {
		return nil, err
	}

	categoryKeyJSON, err := json.Marshal(categoryKey)
	if err != nil {
		return nil, err
	}
	err = ctx.GetStub().PutState("latestKeyAd", categoryKeyJSON)
	if err != nil {
		return nil, err
	}

	adRes := AdRes{
		AdSeq: keyString,
	}

	return &adRes, nil
}

// 동의서 갱신
func (s *SmartContract) UpdateDoc(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, document string, useYn bool) (*AdRes, error) {

	address, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}

	org, err := findOrg(ctx, address)
	if err != nil {
		return nil, err
	}

	// 동의문서 해시화
	userInfohash := sha256.New()
	userInfohash.Write([]byte(document))

	md := userInfohash.Sum(nil)
	docMdStr := hex.EncodeToString(md)

	// 동의문서 찾기
	ad, adSeq, err := findAd(ctx, address, docMdStr)
	if err != nil {
		return nil, err
	}

	adEdit := Ad{
		Address:  org.OrgAddress,
		Document: ad.Document,
		UseYn:    useYn,
	}

	adJSON, err := json.Marshal(adEdit)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(adSeq, adJSON)
	if err != nil {
		return nil, err
	}

	adRes := AdRes{
		AdSeq: adSeq,
	}

	return &adRes, nil
}

// 동의서 조회
func (s *SmartContract) ReadDoc(ctx contractapi.TransactionContextInterface, adSeq string, document string) (*AdReadRes, error) {

	adReadRes := AdReadRes{}

	// 동의문서 해시화
	userInfohash := sha256.New()
	userInfohash.Write([]byte(document))

	md := userInfohash.Sum(nil)
	docMdStr := hex.EncodeToString(md)

	// 조직 seq 쿼리
	adJSON, err := ctx.GetStub().GetState(adSeq)
	if err != nil {
		adReadRes.Result = false
		adReadRes.UseYn = false
		return &adReadRes, err
	}

	ad := Ad{}
	json.Unmarshal(adJSON, &ad)

	if ad.Document != docMdStr {
		adReadRes.Result = false
		adReadRes.UseYn = false
		return &adReadRes, err
	}

	adReadRes.Result = true
	adReadRes.UseYn = ad.UseYn

	return &adReadRes, err

}

// 동의정보 등록
func (s *SmartContract) CreateAgreement(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, adSeq string, useYn bool) (*AgreeRes, error) {

	var categoryKey = CategoryKey{}
	categoryBytes, err := generateKey(ctx, "latestKeyAgree", "AGREE")
	if err != nil {
		return nil, err
	}
	json.Unmarshal(categoryBytes, &categoryKey)
	keyidx := strconv.Itoa(categoryKey.Idx)
	var keyString = categoryKey.Key + keyidx

	userAddress, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}

	// 유저 찾기
	user, err := findUser(ctx, userAddress, "")
	if err != nil {
		return nil, err
	}

	// 조직 seq 쿼리
	adJSON, err := ctx.GetStub().GetState(adSeq)
	if err != nil {
		return nil, err
	}

	ad := Ad{}
	json.Unmarshal(adJSON, &ad)

	if !ad.UseYn {
		return nil, err
	}

	agree := Agree{
		AgreeSeq: keyString,
		Address:  user.Address,
		AdSeq:    adSeq,
		AgreeYn:  useYn,
	}

	agreeJSON, err := json.Marshal(agree)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(keyString, agreeJSON)
	if err != nil {
		return nil, err
	}

	categoryKeyJSON, err := json.Marshal(categoryKey)
	if err != nil {
		return nil, err
	}
	err = ctx.GetStub().PutState("latestKeyAgree", categoryKeyJSON)
	if err != nil {
		return nil, err
	}

	agreeRes := AgreeRes{
		AgreeSeq: keyString,
	}

	return &agreeRes, nil
}

// 동의정보 수정
func (s *SmartContract) UpdateAgreement(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, adSeq string, useYn bool) (*AgreeRes, error) {

	orgAddress, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}

	// 유저 찾기
	agree, agreeSeq, err := findAgree(ctx, orgAddress, adSeq)
	if err != nil {
		return nil, err
	}

	// 동의서 seq 쿼리
	adJSON, err := ctx.GetStub().GetState(adSeq)
	if err != nil {
		return nil, err
	}

	ad := Ad{}
	json.Unmarshal(adJSON, &ad)

	if !ad.UseYn {
		return nil, err
	}

	agreeEdit := Agree{
		Address: agree.Address,
		AdSeq:   adSeq,
		AgreeYn: useYn,
	}

	agreeJSON, err := json.Marshal(agreeEdit)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(agreeSeq, agreeJSON)
	if err != nil {
		return nil, err
	}

	agreeRes := AgreeRes{
		AgreeSeq: agreeSeq,
	}

	return &agreeRes, nil
}

// 동의정보 조회
func (s *SmartContract) ReadAgreement(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, agreeSeq string, userSeq string) (*AgreeReadRes, error) {

	orgAddress, err := getAddress(privateKeyEncoded)

	// 유저 정보 seq 쿼리
	userJSON, err := ctx.GetStub().GetState(userSeq)
	if err != nil {
		return nil, err
	}
	user := User{}
	json.Unmarshal(userJSON, &user)

	// 동의정보 seq 쿼리
	agreeJSON, err := ctx.GetStub().GetState(agreeSeq)
	if err != nil {
		return nil, err
	}
	agree := Agree{}
	json.Unmarshal(agreeJSON, &agree)

	if agree.Address != user.Address {
		return nil, err
	}

	// 동의서 seq 쿼리
	adJSON, err := ctx.GetStub().GetState(agree.AdSeq)
	if err != nil {
		return nil, err
	}

	ad := Ad{}
	json.Unmarshal(adJSON, &ad)

	if ad.Address != orgAddress {
		return nil, err
	}

	agreeReadRes := AgreeReadRes{
		Result:  true,
		AgreeYn: agree.AgreeYn,
	}

	return &agreeReadRes, err

}

//// 의료정보 사용기록

// 의료정보 요청 기록(사용자)
func (s *SmartContract) CreateMediInfo(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, agreeSeq string) (*MediRes, error) {

	var categoryKey = CategoryKey{}
	categoryBytes, err := generateKey(ctx, "latestKeyMedi", "MEDI")
	if err != nil {
		return nil, err
	}
	json.Unmarshal(categoryBytes, &categoryKey)
	keyidx := strconv.Itoa(categoryKey.Idx)
	var keyString = categoryKey.Key + keyidx

	userAddress, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}

	// 유저 찾기
	user, err := findUser(ctx, userAddress, "")
	if err != nil {
		return nil, err
	}

	// 조직 seq 쿼리
	agreeJSON, err := ctx.GetStub().GetState(agreeSeq)
	if err != nil {
		return nil, err
	}

	agree := Agree{}
	json.Unmarshal(agreeJSON, &agree)

	if !agree.AgreeYn {
		return nil, err
	}

	medi := Medi{
		MediSeq:  keyString,
		UserSeq:  user.UserSeq,
		AgreeSeq: agreeSeq,
		Status:   "request",
	}

	mediJSON, err := json.Marshal(medi)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(keyString, mediJSON)
	if err != nil {
		return nil, err
	}

	categoryKeyJSON, err := json.Marshal(categoryKey)
	if err != nil {
		return nil, err
	}
	err = ctx.GetStub().PutState("latestKeyMedi", categoryKeyJSON)
	if err != nil {
		return nil, err
	}

	mediRes := MediRes{
		MediSeq: keyString,
	}

	return &mediRes, nil
}

// 의료정보 요청  확인(기관)
func (s *SmartContract) ConfirmMedi(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, agreeSeq string) ([]Medi, error) {

	orgAddress, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}

	// 유저 찾기
	org, err := findOrg(ctx, orgAddress)
	if err != nil {
		return nil, err
	}

	return findMedi(ctx, org.OrgAddress, agreeSeq)
}

// 의료정보 요청 검색(슈퍼 권한)
func (s *SmartContract) SearchMedi(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, userSeq string, orgSeq string, status string) ([]Medi, error) {

	// 슈퍼 권한 체크
	orgSuperAddress, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}
	orgInfo, err := findOrg(ctx, orgSuperAddress)
	if err != nil {
		return nil, errors.New("org not fount")
	}
	if orgInfo.OrgType != "super" {
		return nil, errors.New("not super user")
	}

	return findMediSearch(ctx, userSeq, orgSeq, status)

}

// 의료정보 요청 결과 등록
func (s *SmartContract) CreateMediDetail(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, mediSeq string, agreeSeq string, status string, mediInfo string) (*MediDetailRes, error) {

	//카테고리 키 생성
	var categoryKey = CategoryKey{}
	categoryBytes, err := generateKey(ctx, "latestKeyMediDetail", "MEDIDETAIL")
	if err != nil {
		return nil, err
	}
	json.Unmarshal(categoryBytes, &categoryKey)
	keyidx := strconv.Itoa(categoryKey.Idx)
	var keyString = categoryKey.Key + keyidx

	orgAddress, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}

	// 조직 존재여부 확인
	_, err = findOrg(ctx, orgAddress)
	if err != nil {
		return nil, err
	}

	// 동의문서 해시화
	userInfohash := sha256.New()
	userInfohash.Write([]byte(mediInfo))

	md := userInfohash.Sum(nil)
	mediInfoMdStr := hex.EncodeToString(md)

	// 동의정보 seq 쿼리
	_, err = ctx.GetStub().GetState(agreeSeq)
	if err != nil {
		return nil, err
	}

	// 의료정보 seq 쿼리
	mediJSON, err := ctx.GetStub().GetState(mediSeq)
	if err != nil {
		return nil, err
	}

	medi := Medi{}
	json.Unmarshal(mediJSON, &medi)

	medi.Status = status

	mediJSON, err = json.Marshal(medi)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(mediSeq, mediJSON)
	if err != nil {
		return nil, err
	}

	mediDetail := MediDetail{
		MediDetailSeq: keyString,
		Address:       orgAddress,
		AgreeSeq:      agreeSeq,
		MediSeq:       mediSeq,
		MediInfo:      mediInfoMdStr,
		Status:        status,
	}

	mediDetailJSON, err := json.Marshal(mediDetail)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(keyString, mediDetailJSON)
	if err != nil {
		return nil, err
	}

	// 카테고리 최신화
	categoryKeyJSON, err := json.Marshal(categoryKey)
	if err != nil {
		return nil, err
	}
	err = ctx.GetStub().PutState("latestKeyMediDetail", categoryKeyJSON)
	if err != nil {
		return nil, err
	}

	mediDetailRes := MediDetailRes{
		MediDetailSeq: keyString,
	}

	return &mediDetailRes, nil
}

// 의료정보 요청 결과 갱신
func (s *SmartContract) UpdateMediDetail(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, mediSeq string, agreeSeq string, mediInfo string) (*MediDetailRes, error) {

	// 슈퍼 권한 체크
	address, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}
	orgInfo, err := findOrg(ctx, address)

	if err != nil {
		return nil, err
	}
	if orgInfo.OrgType != "super" {
		return nil, err
	}

	// 유저 찾기
	mediDetail, err := findMediDetail(ctx, mediSeq, agreeSeq)
	if err != nil {
		return nil, err
	}

	// 의료정보 seq 쿼리
	mediDetailJSON, err := ctx.GetStub().GetState(mediDetail.MediDetailSeq)
	if err != nil {
		return nil, err
	}

	mediDetailEdit := MediDetail{}
	json.Unmarshal(mediDetailJSON, &mediDetailEdit)

	// 동의문서 해시화
	userInfohash := sha256.New()
	userInfohash.Write([]byte(mediInfo))

	md := userInfohash.Sum(nil)
	mediInfoMdStr := hex.EncodeToString(md)

	if mediInfoMdStr == mediDetail.MediInfo {
		mediDetailEdit.Status = "success"
	} else {
		mediDetailEdit.Status = "error"
	}

	mediDetailJSON, err = json.Marshal(mediDetailEdit)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(mediDetail.MediDetailSeq, mediDetailJSON)
	if err != nil {
		return nil, err
	}

	mediDetailRes := MediDetailRes{
		MediDetailSeq: mediDetail.MediDetailSeq,
	}

	return &mediDetailRes, nil
}

// 의료정보 요청 결과 비교
func (s *SmartContract) CompareMediDetail(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, mediSeq string, agreeSeq string, mediInfo string) (*MediCompareRes, error) {

	// super check
	address, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}
	orgInfo, err := findOrg(ctx, address)

	if err != nil {
		return nil, err
	}
	if orgInfo.OrgType != "super" {
		return nil, err
	}

	// find user
	mediDetail, err := findMediDetail(ctx, mediSeq, agreeSeq)
	if err != nil {
		return nil, err
	}

	// document consent hash
	userInfohash := sha256.New()
	userInfohash.Write([]byte(mediInfo))

	md := userInfohash.Sum(nil)
	mediInfoMdStr := hex.EncodeToString(md)

	mediCompareRes := MediCompareRes{
		Result: false,
	}

	if mediInfoMdStr == mediDetail.MediInfo {
		mediCompareRes.Result = true
	}

	return &mediCompareRes, nil
}

// 의료정보 요청 결과 조회
func (s *SmartContract) SearchMediDetail(ctx contractapi.TransactionContextInterface, privateKeyEncoded string, mediDetailSeq string) ([]MediDetail, error) {

	// super check
	orgAddress, err := getAddress(privateKeyEncoded)
	if err != nil {
		return nil, err
	}
	orgInfo, err := findOrg(ctx, orgAddress)

	if err != nil {
		return nil, err
	}
	if orgInfo.OrgType == "super" {
		return findMediDetailSearch(ctx, "", mediDetailSeq)
	} else {
		return findMediDetailSearch(ctx, orgAddress, mediDetailSeq)
	}

}
