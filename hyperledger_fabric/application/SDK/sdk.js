'use strict';
const { FileSystemWallet, Gateway } = require('fabric-network');
var path = require('path');
const ccpPath = path.resolve(__dirname, '..', 'connection.json');
async function send(type, func, walletid, args, res){
    try {
        /** set channelname */
        const channelname = "";
        /** set contractname*/
        const contractname = "";


        const walletPath = path.join(process.cwd(), '..', 'wallet');
        const wallet = new FileSystemWallet(walletPath); 
        console.log(`Wallet path: ${walletPath}`);

        console.log(`Wallet Id : ${walletid} `)

        const userExists = await wallet.exists(walletid);

        if (!userExists) {
            console.log('An identity for the user "'+walletid+'" does not exist in the wallet');
            console.log('Run the registUser.js application before retrying');
            return;
        }
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: walletid, discovery: { enabled: true, asLocalhost: true } });
        const network = await gateway.getNetwork(channelname);
        const contract = network.getContract(contractname);
        if(type){
            await contract.submitTransaction(func, ...args);
            console.log('Transaction has been submitted');
            await gateway.disconnect();
            res.send('success');
        }else{
            const result = await contract.evaluateTransaction(func, ...args);
            console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
            res.send(result.toString());
        }
    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        res.send(`Failed to submit transaction: ${error}`);
    }
}
module.exports = {
    send:send
}