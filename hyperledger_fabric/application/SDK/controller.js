'use strict';
var sdk = require('./sdk.js');
var registUsers = require('./registUsers.js')

module.exports = function(app){
    /**
    * name    : /api/registUsers
    * type	  : function
    * comment : regist user in blockchain network
    * @param  	userid
    * @return
    */
    app.get('/api/registUsers', function (req, res) {
        var userid = req.query.userid;
        let args = [userid];
        registUsers.send(args, res);
    });

    /**
    * name    : /api/getAgreeByWalletId
    * type	  : function
    * comment : search agree
    * @param  	wallertid
    * @return
    */
    app.get('/api/getAgreeByWalletId', function (req, res) {
        var walletid = req.query.walletid;
        let args =[walletid];
        sdk.send(false, 'getAgreeByWalletId', walletid, args, res);
    });

    /**
    * name    : /api/getAgreeByWalletIdAndAgreeKey
    * type	  : function
    * comment : search agree
    * @param  	wallertid
    * @param  	agreeKey
    * @return
    */
    app.get('/api/getAgreeByWalletIdAndAgreeKey', function (req, res) {
        var walletid = req.query.walletid;
        var agreekey = req.query.agreekey;
        let args =[walletid,agreekey];
        sdk.send(false, 'getAgreeByWalletIdAndAgreeKey', walletid, args, res);
    });

    /**
    * name    : /api/getAllAgree
    * type	  : function
    * comment : search agree
    * @param  	wallertid
    * @param  	agreeKey
    * @return
    */
    app.get('/api/getAllAgree', function (req, res) {
        var walletid = req.query.walletid;
        let args =[];
        sdk.send(false, 'getAllAgree', walletid, args, res);
    });

    /**
    * name    : /api/setAgree
    * type	  : function
    * comment : insert agree data
    * @param  	walletid
    * @param  	userid
    * @param  	agreekey
    * @param  	agree
    * @return
    */
    app.get('/api/setAgree', function (req, res) {
        var walletid = req.query.walletid;
        var userid = req.query.userid;
        var agreekey = req.query.agreekey;
        var agree = req.query.agree;
        let args = [userid,agreekey,agree];
        sdk.send(true, 'setAgree',walletid, args, res);
    });

    /**
    * name    : /api/getHistory
    * type	  : function
    * comment : search history
    * @param  	walletid
    * @param  	key
    * @return
    */
    app.get('/api/getHistory', function (req, res) {
        var walletid = req.query.walletid;
        var key = req.query.key;
        let args =[key];
        sdk.send(false, 'getHistory', walletid, args, res);
    });

    /**
    * name    : /api/getHistoryByWalletId
    * type	  : function
    * comment : search history
    * @param  	walletid
    * @return
    */
    app.get('/api/getHistoryByWalletId', function (req, res) {
        var walletid = req.query.walletid;
        let args =[walletid];
        sdk.send(false, 'getAgreeByWalletId', walletid, args, res);
    });

    /**
    * name    : /api/getHistoryByWalletIdAndHistoryType
    * type	  : function
    * comment : search history
    * @param  	walletid
    * @param  	historytype
    * @return
    */
    app.get('/api/getHistoryByWalletIdAndHistoryType', function (req, res) {
        var walletid = req.query.walletid;
        var historytype = req.query.historytype;
        let args =[walletid,historytype];
        sdk.send(false, 'getHistoryByWalletIdAndHistoryType', walletid, args, res);
    });
    
    /**
    * name    : /api/getAllHistory
    * type	  : function
    * comment : search all history
    * @param  	walletid
    * @return
    */
    app.get('/api/getAllHistory', function (req, res) {
        var walletid = req.query.walletid;
        let args =[];
        sdk.send(false, 'getAllAgree', walletid, args, res);
    });

    /**
    * name    : /api/setHistory
    * type	  : function
    * comment : insert history
    * @param  	walletid
    * @param  	userid
    * @param  	historytype
    * @return
    */
    app.get('/api/setHistory', function (req, res) {
        var walletid = req.query.walletid;
        var userid = req.query.userid;
        var historytype = req.query.historytype;
        let args = [userid,historytype];
        sdk.send(true, 'setHistory',walletid, args, res);
    });

    /**
    * name    : /api/updateHistory_Run
    * type	  : function
    * comment : update hisyory Run
    * @param  	walletid
    * @param  	userid
    * @param  	historytype
    * @param  	runtype
    * @return
    */
    app.get('/api/updateHistory_Run', function (req, res) {
        var walletid = req.query.walletid;
        var userid = req.query.userid;
        var historytype = req.query.historytype;
        var runtype = req.query.runtype;
        let args = [userid,historytype,runtype];
        sdk.send(true, 'updateHistory_Run',walletid, args, res);
    });
}
