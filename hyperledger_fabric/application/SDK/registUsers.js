'use strict';

const { FileSystemWallet, Gateway, X509WalletMixin } = require('fabric-network');
const path = require('path');

const ccpPath = path.resolve(__dirname, '..', 'connection.json');

async function create(args, res) {
    try {
        /** set adminName */
        const adminName = ""; //ex.admin
        /** set affiliationName*/
        const affiliationName = ""; //ex.org1.department1
        /** set roleName*/
        const roleName = ""; //ex.client
        /** set orgName*/
        const orgName = ""; //ex.Sales1Org

        const walletPath = path.join(process.cwd(), '..', 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        var userid = args[0];

        const userExists = await wallet.exists(userid);
        if (userExists) {
            console.log('An identity for the user "'+userid+'"already exists in the wallet');
            return;
        }

        const adminExists = await wallet.exists(adminName);
        if (!adminExists) {
            console.log('An identity for the admin user "admin" does not exist in the wallet');
            console.log('Run the enrollAdmin.js application before retrying');
            return;
        }

        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: adminName, discovery: { enabled: true, asLocalhost: true } });

        const ca = gateway.getClient().getCertificateAuthority();
        const adminIdentity = gateway.getCurrentIdentity();

        const secret = await ca.register({ affiliation: affiliationName, enrollmentID: userid, role: roleName }, adminIdentity);
        const enrollment = await ca.enroll({ enrollmentID: userid, enrollmentSecret: secret });
        const userIdentity = X509WalletMixin.createIdentity(orgName, enrollment.certificate, enrollment.key.toBytes());
        await wallet.import(userid, userIdentity);
        console.log('Successfully registered and enrolled admin user "'+userid+'" and imported it into the wallet');
    
        res.send('success');

    } catch (error) {
        console.error(`Failed to register user "${userid}":${error}`);
        process.exit(1);
    }
}

module.exports = {
    send:create
}