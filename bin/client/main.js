'use strict';

var game = null
var gParam = {ws_server: "127.0.0.1:8989", user_name: "", joinroom: null, platform: "PC", app_token: null}
//var gParam = {ws_server:"172.24.222.54:8989/ws", user_name:"", joinroom:null, platform:"PC", app_token:null}


var gImageDir = "assets/2x/"
var gFontScale = 1.0;
var Native

function bindNative() {
    if (PLAT == "IOS") {
        Native = {}

        enableNativeLog();

        try {


            connectWebViewJavascriptBridge(function (bridge) {
                bridge.init(function (message, responseCallback) {
                    console.log("WebViewJavascriptBridge init OK");
                });
                bridge.registerHandler("quitApp", function () {
                    game.Native.yesOrNoPopupWindow("退出游戏", "你确定要退出游戏吗？", "取消", "确定", function (data) {
                        if (data.sender == "popButton2") {
                            game.Native.quitToApp();
                        }

                    });

                });

                bridge.callHandler("getNativeConfig", null, function (data) {

                    console.log(data.ws_server);
                    data.platform = PLAT
                    startGame(data);
                });

                Native.quitToApp = function () {
                    game.betApi.leaveRoom();
                    bridge.callHandler("quitToApp");
                    console.log("quitToApp")
                };

                Native.showProfile = function (userid) {
                    bridge.callHandler("showProfile", {userid: userid});
                    console.log("showProfile");
                },

                    Native.confrimPopupWindow = function (title, text, buttonText, callback) {
                        if (buttonText == undefined || buttonText == null) {
                            buttonText = "确定";
                        }
                        bridge.callHandler("confrimPopupWindow", {
                            pop_title: title,
                            pop_text: text,
                            pop_btn1_text: buttonText
                        }, function (data) {
                            callback(data);
                        });
                        console.log("confrimPopupWindow")
                    };

                Native.yesOrNoPopupWindow = function (title, text, button1Text, button2Text, callback) {
                    if (button1Text == undefined || button1Text == null) {
                        button1Text = "取消";
                    }

                    if (button2Text == undefined || button2Text == null) {
                        button2Text = "确定";
                    }

                    bridge.callHandler("yesOrNoPopupWindow", {
                        pop_title: title,
                        pop_text: text,
                        pop_btn1_text: button1Text,
                        pop_btn2_text: button2Text
                    }, function (data) {
                        callback(data);
                    });
                    console.log("yesOrNoPopupWindow")
                };
            });


        } catch (e) {
            console.log(e);
        }


    } else {
        startGame({platform: PLAT});
    }

}

function enableNativeLog() {
    console = new Object();
    console.log = function (log, other) {
        var iframe = document.createElement("IFRAME");
        var otherstring = ""
        if (other != undefined) {
            otherstring = other;
        }

        iframe.setAttribute("src", "ios-log:#iOS#" + log + " " + otherstring);
        document.documentElement.appendChild(iframe);
        iframe.parentNode.removeChild(iframe);
        iframe = null;
    };
    console.debug = console.log;
    console.info = console.log;
    console.warn = console.log;
    console.error = console.log;
}


function connectWebViewJavascriptBridge(callback) {
    if (window.WebViewJavascriptBridge) {
        callback(WebViewJavascriptBridge)
    } else {
        document.addEventListener('WebViewJavascriptBridgeReady', function () {
            callback(WebViewJavascriptBridge)
        }, false)
    }
}


function _fontString(size, fontname) {
    if (fontname == undefined) {
        //fontname = "Impact"
        fontname = "Apple LiSung Light"
    }
    ;

    return (size * gFontScale) + "px " + fontname
}

function startGame(gameParam) {
    try {
        // merge property

        for (var p in gParam) {
            var value = gameParam[p];
            if (value != undefined && value != null) {
                gParam[p] = gameParam[p];
            }
        }
        if (gParam["platform"] == "IOS") {
            var deviceWidth = 1136;
            var deviceHeight = 640;
            document.body.setAttribute("orient", "landscape");
            gImageDir = "assets/1x/"
            gFontScale = 0.6;
            // 使用 1x的图片效果更差，算了，直接用2x的

            gImageDir = "assets/2x/"
            gFontScale = 1.2;

            game = new Phaser.Game(deviceWidth, deviceHeight, Phaser.CANVAS, "gamediv");
        } else {
            game = new Phaser.Game("100", "100", Phaser.CANVAS, "gamediv");
        }

        game.Native = Native;

        game.betApi = new BetApi();
        if (gParam["app_token"] == undefined || gParam["app_token"] == null) {
            game.state.add("LoginState", LoginState);
        }

        game.state.add("MainState", MainState);
        //gParam["app_token"] = "testUSer"

        if (gParam["app_token"] != undefined && gParam["app_token"] != null) {
            game.state.start("MainState");
        } else {
            //gParam["app_token"] = "testUSer"
            game.state.start("LoginState");
            //game.state.start("MainState");
        }
    } catch (e) {
        console.log("error ! ", e);
    }
}

function gameQuit(cause) {
    if (gParam["platform"] == "IOS") {
        //game.Native.quitToApp(cause);

        game.Native.confrimPopupWindow("你的钱输光了！！", "你的积分为0， 即将被踢出游戏", "确认", function (data) {
            console.log("Good lllll");
        });


        /*
         game.Native.yesOrNoPopupWindow("你的钱输光了！！","你的积分为0， 即将被踢出游戏", "放弃", "确认",function(data){
         console.log("you click:", data.sender);
         });
         */
    } else {
        game.state.states["MainState"].actionOnExit();
    }
}
