System.register(['angular2/core', 'angular2/router', 'angular2/http', 'rxjs/add/operator/map', 'rxjs/add/operator/catch', './ng2-qrcode.js'], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
        var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
        if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
        else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
        return c > 3 && r && Object.defineProperty(target, key, r), r;
    };
    var __metadata = (this && this.__metadata) || function (k, v) {
        if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
    };
    var core_1, router_1, http_1, http_2, ng2_qrcode_ts_1;
    var loadWalletComponent, DisplayModeEnum;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (router_1_1) {
                router_1 = router_1_1;
            },
            function (http_1_1) {
                http_1 = http_1_1;
                http_2 = http_1_1;
            },
            function (_1) {},
            function (_2) {},
            function (ng2_qrcode_ts_1_1) {
                ng2_qrcode_ts_1 = ng2_qrcode_ts_1_1;
            }],
        execute: function() {
            let loadWalletComponent = class loadWalletComponent {
                //Constructor method for load HTTP object
                constructor(http) {
                    this.http = http;
                    this.displayModeEnum = DisplayModeEnum;
                }
                //Init function for load default value
                ngOnInit() {
                    this.displayMode = DisplayModeEnum.first;
                    this.loadWallet();
                    this.loadProgress();
                    //Set interval function for load wallet every 15 seconds
                    setInterval(() => {
                        this.loadWallet();
                        console.log("Refreshing balance");
                    }, 15000);
                    //Enable Send tab "textbox" and "Ready" button by default
                    this.sendDisable = true;
                    this.readyDisable = false;
                    this.pendingTable = [];
                    if (localStorage.getItem('historyAddresses') != null) {
                        this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
                    }
                    else {
                        localStorage.setItem('historyAddresses', JSON.stringify([]));
                        this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
                    }
                }
                //Ready button function for disable "textbox" and enable "Send" button for ready to send coin
                ready(spendId, spendaddress, spendamount) {
                    if (!spendId) {
                        alert("Please select from id");
                        return false;
                    }
                    if (!spendaddress) {
                        alert("Please enter pay to");
                        return false;
                    }
                    if (!spendamount) {
                        alert("Please enter amount");
                        return false;
                    }
                    this.readyDisable = true;
                    this.sendDisable = false;
                }
                //Load wallet function
                loadWallet() {
                    this.http.post('/wallets', '')
                        .map((res) => res.json())
                        .subscribe(data => {
                        this.wallets = data;
                        console.log('walets data is : ' + this.wallets);
                        //Load Balance for each wallet
                        //Set http headers
                        var headers = new http_2.Headers();
                        headers.append('Content-Type', 'application/x-www-form-urlencoded');
                        var inc = 0;
                        var myInc = 0;
                        for (var item in data) {
                            var wallet_id = data[inc].meta.filename;
                            //var wallet_id = item.meta.filename;
                            console.log('wallet_id is : ' + wallet_id);
                            console.log('inc is : ' + inc);
                            //Post method executed
                            //this.http.post('/wallet/balance', JSON.stringify({id: address}), {headers: headers})
                            this.http.get('/wallet/balance?id=' + wallet_id, { headers: headers })
                                .map((res) => res.json())
                                .subscribe(
                            //Response from API
                            response => {
                                console.log('load done: ' + myInc);
                                this.wallets[myInc].balance = response.confirmed.coins / 1000000;
                                myInc++;
                            }, err => console.log("Error on load balance: " + err), () => console.log('Balance load done'));
                            inc++;
                        }
                        //Load Balance for each wallet end
                    }, err => console.log("Error on load wallet: " + err), () => console.log('Wallet load done'));
                }
                //Load progress function for Skycoin
                loadProgress() {
                    //Post method executed
                    this.http.post('/blockchain/progress', '')
                        .map((res) => res.json())
                        .subscribe(
                    //Response from API
                    response => { this.progress = (parseInt(response.current, 10) + 1) / parseInt(response.highest, 10) * 100; }, err => console.log("Error on load progress: " + err), () => console.log('Progress load done:' + this.progress));
                }
                //Switch tab function
                switchTab(mode, wallet) {
                    //"Textbox" and "Ready" button enable in Send tab while switching tabs
                    this.sendDisable = true;
                    this.readyDisable = false;
                    this.displayMode = mode;
                    if (wallet) {
                        this.spendid = wallet.meta.filename;
                    }
                }
                //Show QR code function for show QR popup
                showQR(wallet) {
                    this.QrAddress = wallet.entries[0].address;
                    this.QrIsVisible = true;
                }
                //Hide QR code function for hide QR popup
                hideQrPopup() {
                    this.QrIsVisible = false;
                }
                //Show wallet function for view New wallet popup
                showNewWalletDialog() {
                    this.NewWalletIsVisible = true;
                }
                //Hide wallet function for hide New wallet popup
                hideWalletPopup() {
                    this.NewWalletIsVisible = false;
                }
                //Add new wallet function for generate new wallet in Skycoin
                createNewWallet() {
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    //Post method executed
                    this.http.post('/wallet/create', JSON.stringify({ name: '' }), { headers: headers })
                        .map((res) => res.json())
                        .subscribe(response => {
                        //Hide new wallet popup
                        this.NewWalletIsVisible = false;
                        alert("New wallet created successfully");
                        //Load wallet for refresh list
                        this.loadWallet();
                    }, err => console.log("Error on create new wallet: " + err), () => console.log('New wallet create done'));
                }
                //Edit existing wallet function
                editWallet(wallet) {
                    this.EditWalletIsVisible = true;
                    this.walletId = wallet.meta.filename;
                }
                //Hide edit wallet function
                hideEditWalletPopup() {
                    this.EditWalletIsVisible = false;
                }
                //Update wallet function for update wallet label
                updateWallet(walletid, walletName) {
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    var stringConvert = 'name=' + walletName + '&id=' + walletid;
                    //Post method executed
                    this.http.post('/wallet/update', stringConvert, { headers: headers })
                        .map((res) => res.json())
                        .subscribe(response => {
                        //Hide new wallet popup
                        this.EditWalletIsVisible = false;
                        alert("Wallet updated successfully");
                        //Load wallet for refresh list
                        this.loadWallet();
                    }, err => console.log("Error on update wallet: " + JSON.stringify(err)), () => console.log('Update wallet done'));
                }
                //Load wallet seed function
                openLoadWallet(walletName, seed) {
                    this.loadSeedIsVisible = true;
                }
                //Hide load wallet seed function
                hideLoadSeedWalletPopup() {
                    this.loadSeedIsVisible = false;
                }
                //Load wallet seed function for create new wallet with name and seed
                createWalletSeed(walletName, seed) {
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    var stringConvert = 'name=' + walletName + '&seed=' + seed;
                    //Post method executed
                    this.http.post('/wallet/create', stringConvert, { headers: headers })
                        .map((res) => res.json())
                        .subscribe(response => {
                        //Hide load wallet seed popup
                        this.loadSeedIsVisible = false;
                        //Load wallet for refresh list
                        this.loadWallet();
                    }, err => console.log("Error on create load wallet seed: " + JSON.stringify(err)), () => console.log('Load wallet seed done'));
                }
                spend(spendid, spendaddress, spendamount) {
                    //Set local storage for history
                    if (localStorage.getItem('historyTable') != null) {
                        this.historyTable = JSON.parse(localStorage.getItem('historyTable'));
                    }
                    else {
                        localStorage.setItem('historyTable', JSON.stringify([]));
                        this.historyTable = JSON.parse(localStorage.getItem('historyTable'));
                    }
                    this.historyTable.push({ address: spendaddress, amount: spendamount });
                    localStorage.setItem('historyTable', JSON.stringify(this.historyTable));
                    //Set local storage for addresses history
                    if (localStorage.getItem('historyAddresses') != null) {
                        this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
                    }
                    else {
                        localStorage.setItem('historyAddresses', JSON.stringify([]));
                        this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
                    }
                    this.addresses.push({ address: spendaddress, amount: spendamount });
                    localStorage.setItem('historyAddresses', JSON.stringify(this.addresses));
                    this.readyDisable = true;
                    this.sendDisable = true;
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    var stringConvert = 'id=' + spendid + '&coins=' + spendamount * 1000000 + "&fee=1&hours=1&dst=" + spendaddress;
                    //Post method executed
                    this.http.post('/wallet/spend', stringConvert, { headers: headers })
                        .map((res) => res.json())
                        .subscribe(response => {
                        this.pendingTable.push({ complete: 'Completed', address: spendaddress, amount: spendamount });
                        //Load wallet for refresh list
                        this.loadWallet();
                    }, err => {
                        alert(err._body);
                        this.readyDisable = false;
                        this.sendDisable = true;
                        this.pendingTable.push({ complete: 'Pending', address: spendaddress, amount: spendamount });
                    }, () => console.log('Spend successfully'));
                }
            };
            loadWalletComponent = __decorate([
                core_1.Component({
                    selector: 'load-wallet',
                    directives: [router_1.ROUTER_DIRECTIVES, ng2_qrcode_ts_1.QRCodeComponent],
                    providers: [],
                    templateUrl: 'app/templates/wallet.html'
                }), 
                __metadata('design:paramtypes', [http_1.Http])
            ], loadWalletComponent);
            exports_1("loadWalletComponent", loadWalletComponent);
            //Set default enum value for tabs
            (function (DisplayModeEnum) {
                DisplayModeEnum[DisplayModeEnum["first"] = 0] = "first";
                DisplayModeEnum[DisplayModeEnum["second"] = 1] = "second";
                DisplayModeEnum[DisplayModeEnum["third"] = 2] = "third";
            })(DisplayModeEnum || (DisplayModeEnum = {}));
        }
    }
});

//# sourceMappingURL=app.loadWallet.js.map
