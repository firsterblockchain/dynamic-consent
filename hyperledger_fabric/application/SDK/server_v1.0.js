const express = require('express');
const bodyParser = require('body-parser');
const EthCrypto = require('eth-crypto');
const app = express();

var path = require('path');
var sdk = require('./sdk');

const PORT = 80;
const HOST = '0.0.0.0';

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }))

// 신원 인증
app.post('/auth/user', function (req, res) {
    var key = req.body.key;
    var userName = req.body.userName;
    var userBirthday = req.body.userBirthday;
    var userSex = req.body.userSex;

    var privateKey = EthCrypto.createIdentity().privateKey;

    let args = [privateKey, key, userName, userBirthday, userSex];

    sdk.sendWithPrivatekey(false, 'CreateUser', args, res);
});

app.get('/auth/user', function (req, res) {
    var key = req.query.key;
    var userName = req.query.userName;
    var userBirthday = req.query.userBirthday;
    var userSex = req.query.userSex;

    let args = [key, userName, userBirthday, userSex];
    sdk.send(true, 'ReadUser', args, res);
});

app.post('/auth/user/reissue', function (req, res) {
    var key = req.body.key;
    var userName = req.body.userName;
    var userBirthday = req.body.userBirthday;
    var userSex = req.body.userSex;

    var privateKey = EthCrypto.createIdentity().privateKey;

    let args = [privateKey, key, userName, userBirthday, userSex];

    sdk.sendWithPrivatekey(false, 'ReissueUser', args, res);
});

app.post('/auth/org', function (req, res) {
    var key = req.body.key;
    var orgName = req.body.orgName;
    var orgType = req.body.orgType;

    var privateKey = EthCrypto.createIdentity().privateKey;

    let args = [privateKey, key, orgName, orgType];

    sdk.sendWithPrivatekey(false, 'CreateOrg', args, res);
});

app.post('/auth/reissue', function (req, res) {
    var key = req.body.key;
    var orgName = req.body.orgName;
    var orgType = req.body.orgType;

    var privateKey = EthCrypto.createIdentity().privateKey;

    let args = [privateKey, key, orgName, orgType];

    sdk.sendWithPrivatekey(false, 'ReissueOrg', args, res);
});

app.post('/auth/org/reg', function (req, res) {
    var privateKey = req.body.privateKey;
    var orgSeq = req.body.orgSeq;
    var orgID = req.body.orgID;
    var userName = req.body.userName;
    var userBirthday = req.body.userBirthday;
    var userSex = req.body.userSex;

    let args = [privateKey, orgSeq, orgID, userName, userBirthday, userSex];

    sdk.send(false, 'CreateOrgUser', args, res);
});

app.get('/auth/org', function (req, res) {
    var orgSeq = req.query.orgSeq;
    var orgID = req.query.orgID;
    var userName = req.query.userName;
    var userBirthday = req.query.userBirthday;
    var userSex = req.query.userSex;

    let args = [orgSeq, orgID, userName, userBirthday, userSex];
    sdk.send(true, 'ReadOrgUser', args, res);
});

app.post('/auth/user/lock', function (req, res) {
    var privateKey = req.body.privateKey;
    var userSeq = req.body.userSeq;

    let args = [privateKey, userSeq];

    sdk.send(false, 'LockUser', args, res);
});

app.post('/auth/org/lock', function (req, res) {
    var privateKey = req.body.privateKey;
    var orgSeq = req.body.orgSeq;

    let args = [privateKey, orgSeq];

    sdk.send(false, 'LockOrg', args, res);
});

// 동의서
app.post('/ad/reg', function (req, res) {
    var privateKey = req.body.privateKey;
    var document = req.body.document;
    var useYn = req.body.useYn;

    let args = [privateKey, document, useYn];

    sdk.send(false, 'CreateDoc', args, res);
});

app.post('/ad/edit', function (req, res) {
    var privateKey = req.body.privateKey;
    var document = req.body.document;
    var useYn = req.body.useYn;

    let args = [privateKey, document, useYn];

    sdk.send(false, 'UpdateDoc', args, res);
});

app.get('/ad', function (req, res) {
    var adSeq = req.query.adSeq;
    var document = req.query.document;

    let args = [adSeq, document];
    sdk.send(true, 'ReadDoc', args, res);
});

app.post('/ad/agree/reg', function (req, res) {
    var privateKey = req.body.privateKey;
    var adSeq = req.body.adSeq;
    var agreeYn = req.body.agreeYn;

    let args = [privateKey, adSeq, agreeYn];

    sdk.send(false, 'CreateAgreement', args, res);
});

app.post('/ad/agree/edit', function (req, res) {
    var privateKey = req.body.privateKey;
    var adSeq = req.body.adSeq;
    var agreeYn = req.body.agreeYn;

    let args = [privateKey, adSeq, agreeYn];

    sdk.send(false, 'UpdateAgreement', args, res);
});

app.post('/ad/agree', function (req, res) {
    var privateKey = req.body.privateKey;
    var agreeSeq = req.body.agreeSeq;
    var userSeq = req.body.userSeq;

    let args = [privateKey, agreeSeq, userSeq];

    sdk.send(false, 'ReadAgreement', args, res);
});

// 의료정보 사용기록
app.post('/medi/reg', function (req, res) {
    var privateKey = req.body.privateKey;
    var agreeSeq = req.body.agreeSeq;

    let args = [privateKey, agreeSeq];

    sdk.send(false, 'CreateMediInfo', args, res);
});

app.post('/medi', function (req, res) {
    var privateKey = req.body.privateKey;
    var agreeSeq = req.body.agreeSeq;

    let args = [privateKey, agreeSeq];

    sdk.send(false, 'ConfirmMedi', args, res);
});

app.post('/medi/search', function (req, res) {
    var privateKey = req.body.privateKey;
    var userSeq = req.body.userSeq;
    var orgSeq = req.body.orgSeq;
    var status = req.body.status;

    let args = [privateKey, userSeq, orgSeq, status];

    sdk.send(false, 'SearchMedi', args, res);
});

app.post('/medi/detail/reg', function (req, res) {
    var privateKey = req.body.privateKey;
    var mediSeq = req.body.mediSeq;
    var agreeSeq = req.body.agreeSeq;
    var status = req.body.status;
    var mediInfo = req.body.mediInfo;

    let args = [privateKey, mediSeq, agreeSeq, status, mediInfo];

    sdk.send(false, 'CreateMediDetail', args, res);
});

app.post('/medi/detail/edit', function (req, res) {
    var privateKey = req.body.privateKey;
    var mediSeq = req.body.mediSeq;
    var agreeSeq = req.body.agreeSeq;
    var mediInfo = req.body.mediInfo;

    let args = [privateKey, mediSeq, agreeSeq, mediInfo];

    sdk.send(false, 'UpdateMediDetail', args, res);
});

app.post('/medi/detail/compare', function (req, res) {
    var privateKey = req.body.privateKey;
    var mediSeq = req.body.mediSeq;
    var agreeSeq = req.body.agreeSeq;
    var mediInfo = req.body.mediInfo;

    let args = [privateKey, mediSeq, agreeSeq, mediInfo];

    sdk.send(false, 'CompareMediDetail', args, res);
});

app.post('/medi/detail/search', function (req, res) {
    var privateKey = req.body.privateKey;
    var mediDetailSeq = req.body.mediDetailSeq;

    let args = [privateKey, mediDetailSeq];

    sdk.send(false, 'SearchMediDetail', args, res);
});

app.use(express.static(path.join(__dirname, './client')));

app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
