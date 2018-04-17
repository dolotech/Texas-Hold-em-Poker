'use strict';

var strVersion = "1.0";
var userName = "cmdTest";
var userID = "";
var loginCertification = false;

function getCookie(name)
{
    var arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");

    if(arr = document.cookie.match(reg))
        return (decodeURI(arr[2]));
    else
        return null;
}

function setCookie(name, value)
{
    var Days = 30;
    var exp = new Date();
    exp.setTime(exp.getTime() + Days*24*60*60*1000);
    var str = name + "="+ encodeURI(value) + "; expires=" + exp.toGMTString();
    document.cookie = str;
}

var callbackOpen = function(data)
{
    console.log("callbackOpen " + data);

    game.betApi.checkVersion(strVersion, function(isOK){
        console.log("checkVersion " + isOK);
    });
};

var callbackClose = function(data)
{
    console.log("callbackClose " + data);
    loginCertification = false;

    game.state.states["MainState"]._disconnectReset();
};

var callbackMessage = function(data)
{
    console.log("callbackMessage " + data);
    if(data.Version && data.Version.Version == strVersion) // checkVersion result
    {
        game.state.states["LoginState"].initUserName();
    }
    else if(!loginCertification) // loginCertification result
    {
        if(data.id)
        {
            userID = data.id;
            game.state.states["MainState"].userID = userID
            game.betApi.setUserID(userID);
            loginCertification = true;

            var LoginState = game.state.states["LoginState"];
            LoginState.hanldeInitUserName();
        }
    }
    else if(data.type == "iq")
    {
        if(data.class == "room")       //查询游戏房间列表
        {
            game.state.states["LoginState"].handleCreateRoom(data);
        }
        else if(data.class == "roomlist")       //查询游戏房间列表
        {
            game.state.states["LoginState"].handleGetRoomList(data);
        }
        else if(data.class == "room")  //查询游戏房间信息
        {

        }
        else if(data.class == "occupant")  //查询玩家信息
        {

        }
    }
    else if(data.type == "message")
    {
    }
    else if(data.type == "presence")
    {
        if(data.action == "active")         //服务器广播进入房间的玩家
        {
        }
        else if(data.action == "gone")      //服务器广播离开房间的玩家
        {
            game.state.states["MainState"].handleGone(data);
        }
        else if(data.action == "join")      //服务器通报加入游戏的玩家
        {
            game.state.states["MainState"].handleJoin(data);
        }
        else if(data.action == "button")    //服务器通报本局庄家
        {
            game.state.states["MainState"].handleButton(data);
        }
        else if(data.action == "preflop")   //服务器通报发牌
        {
            game.state.states["MainState"].handlePreflop(data);
        }
        else if(data.action == "flop")   //发牌
        {
            game.state.states["MainState"].handleFlop(data);
        }
        else if(data.action == "turn")   //发牌
        {
            game.state.states["MainState"].handleTurn(data);
        }
        else if(data.action == "river")   //发牌
        {
            game.state.states["MainState"].handleRiver(data);
        }
        else if(data.action == "pot")       //服务器通报奖池
        {
            game.state.states["MainState"].handlePot(data);
        }
        else if(data.action == "action")    //服务器通报当前下注玩家
        {
            game.state.states["MainState"].handleAction(data);

        }
        else if(data.action == "bet")       //服务器通报玩家下注结果
        {
            game.state.states["MainState"].handleBet(data);

        }
        else if(data.action == "showdown")  //服务器通报摊牌和比牌
        {
            game.state.states["MainState"].handleShowDown(data);
        }
        else if(data.action == "state")  //服务器通报房间信息
        {
            game.state.states["MainState"].handleState(data);
        }
    }
};

var callbackError = function(data)
{
    console.log("callbackError" + data);
};

var LoginState = function() {

    this.scale;                         //全局缩放比
    this.group;
    this.currentPage = 0;
    this.CountPerPage = 10;
    this.currentSelectRoomID = -1;
    this.roomInfoList;

    this.roomList;
    this.roomTextList;
    this.btnPrev;
    this.btnNext;
    this.leName;
    this.leRoomName;
    this.leChangeName;
    this.btnChangeName;
    this.btnLogin;
    this.btnCreate;

    this.textPrev;
    this.tectNext;
    this.textLogin;
    this.tectCreate;
};

LoginState.prototype = {

    preload: function () {
        game.load.image("gamecenterbackground", gImageDir+'background.png');
        game.load.image('buttonblue', gImageDir+'btn-blue.png');
        game.load.image('buttongrey', gImageDir+'btn-grey.png');
        game.load.image('buttonyellow', gImageDir+'btn-yellow.png');
    },

    create: function () {

        game.betApi.connect();
        game.betApi.registerCallback(callbackOpen, callbackClose, callbackMessage, callbackError);

        var imageBK = game.add.image(0, 0, "gamecenterbackground");
        imageBK.visible = false;
        var xScale = game.width / imageBK.width;
        var yScale = game.height / imageBK.height;
        this.scale = xScale < yScale ? xScale : yScale;
        loginCertification = false;

        this.roomInfoList = [];
        this.roomList = [];
        this.roomTextList = [];
        this.group = game.add.group();
        var listWidth = game.width / 2;
        var listItemHeight = game.height / 12;
        var style = { font: _fontString(16), fill: "#0069B2", wordWrapWidth: listWidth * 0.9, align: "left"};
        for(var i = 0; i < 10; i++)
        {
            var listItem = game.add.button(0, i * listItemHeight, 'buttonblue', this.selectRoom);
            listItem.width = listWidth;
            listItem.height = listItemHeight;
            listItem.roomID = "";
            listItem.visible = false;
            this.roomList.push(listItem);
            var roomInfo = game.add.text(listItem.x + 0.05 * listItem.width, listItem.y + 0.45 * listItem.height, "", style);
            roomInfo.scale.setTo(this.scale);
            roomInfo.anchor.set(0, 0.5);
            roomInfo.visible = false;
            this.roomTextList.push(roomInfo);
            this.group.add(listItem);
            this.group.add(roomInfo);
        }

        this.btnPrev = game.add.button(0, 10 * listItemHeight, 'buttonyellow', this.clickPrev, this);
        this.btnPrev.width = listWidth / 2;
        this.btnPrev.height = listItemHeight;
        this.btnNext = game.add.button(listWidth / 2, 10 * listItemHeight, 'buttonyellow', this.clickNext, this);
        this.btnNext.width = listWidth / 2;
        this.btnNext.height = listItemHeight;
        style = { font: _fontString(28), fill: "#CE8D00"};
        this.textPrev = game.add.text(this.btnPrev.x + 0.5 * this.btnPrev.width, this.btnPrev.y + 0.5 * this.btnPrev.height, "上一页", style);
        this.textPrev.anchor.set(0.5);
        this.textPrev.scale.setTo(this.scale);
        this.tectNext = game.add.text(this.btnNext.x + 0.5 * this.btnNext.width, this.btnNext.y + 0.5 * this.btnNext.height, "下一页", style);
        this.tectNext.anchor.set(0.5);
        this.tectNext.scale.setTo(this.scale);
        this.group.add(this.btnPrev);
        this.group.add(this.btnNext);
        this.group.add(this.textPrev);
        this.group.add(this.tectNext);

        this.leName = game.add.text(listWidth * 1.2, game.height * 0.1, "用户名: " + userName, style);
        this.leName.scale.setTo(this.scale);
        this.leRoomName = game.add.text(listWidth * 1.2, game.height * 0.3, "房间名: ", style);
        this.leRoomName.scale.setTo(this.scale);
        this.group.add(this.leName);
        this.group.add(this.leRoomName);

        this.btnChangeName = game.add.button(listWidth * 1.2, game.height * 0.2, 'buttonyellow', this.clickRename, this);
        this.btnChangeName.width = listWidth / 2;
        this.btnChangeName.height = listItemHeight;
        this.leChangeName = game.add.text(this.btnChangeName.x + this.btnChangeName.width / 2, this.btnChangeName.y + this.btnChangeName.height / 2, "重命名", style);
        this.leChangeName.anchor.setTo(0.5);
        this.leChangeName.scale.setTo(this.scale);
        this.group.add(this.btnChangeName);
        this.group.add(this.leChangeName);

        this.btnLogin = game.add.button(listWidth * 1.2, game.height * 0.4, 'buttonyellow', this.clickLogin, this);
        this.btnLogin.width = listWidth / 2;
        this.btnLogin.height = listItemHeight;
        this.btnCreate = game.add.button(listWidth * 1.2, game.height * 0.5, 'buttonyellow', this.clickCreate, this);
        this.btnCreate.width = listWidth / 2;
        this.btnCreate.height = listItemHeight;
        this.textLogin = game.add.text(this.btnLogin.x + 0.5 * this.btnLogin.width, this.btnLogin.y + 0.5 * this.btnLogin.height, "进入房间", style);
        this.textLogin.anchor.set(0.5);
        this.textLogin.scale.setTo(this.scale);
        this.tectCreate = game.add.text(this.btnCreate.x + 0.5 * this.btnCreate.width, this.btnCreate.y + 0.5 * this.btnCreate.height, "建立房间", style);
        this.tectCreate.anchor.set(0.5);
        this.tectCreate.scale.setTo(this.scale);
        this.group.add(this.btnLogin);
        this.group.add(this.btnCreate);
        this.group.add(this.textLogin);
        this.group.add(this.tectCreate);
        this.group.visible = false;
    },

    initUserName:function()
    {
        var name = ""
        
        if (gParam.user_name != null && gParam.user_name != "") {
            name = gParam.user_name;
        } else {
            name = getCookie("name");
        }
        
        if(!name || name.length == 0)
        {
            name = prompt("请输入您的名字","");
        }
        if(name)
        {
            userName = name;
            setCookie("name", userName);
            game.betApi.loginCertification(userName, function(isOK){
                console.log("loginCertification is " +  isOK);
                if(!isOK)
                {
                    userName = "";
                }
            });
        }
        else
        {
            this.initUserName();
        }
    },

    hanldeInitUserName:function()
    {
        var text = "用户名: " + userName;
        this.leName.setText(text);
        this.group.visible = true;
        game.betApi.getRoomList();
    },

    handleCreateRoom:function(data)
    {
        this.currentSelectRoomID = data.room.id;
        this.clickLogin();
    },

    handleGetRoomList:function(data)
    {
        this.roomInfoList = data.rooms;
        if(!this.roomInfoList)
        {
            this.roomInfoList = [];
        }

        while(this.CountPerPage * this.currentPage > this.roomInfoList.length)
        {
            this.currentPage--;
        }

        for(var i = 0; i < this.CountPerPage; i++)
        {
            if(i + this.CountPerPage * this.currentPage < this.roomInfoList.length)
            {
                this.roomList[i].visible = true;
                this.roomTextList[i].visible = true;
                var index = i + this.CountPerPage * this.currentPage;
                this.roomList[i].roomID = this.roomInfoList[index].id;
                var strTitle = "房间ID: " + this.roomInfoList[index].id + "\n小盲注: " + this.roomInfoList[index].sb + "; 大盲注: " + this.roomInfoList[index].bb + "; 当前人数: " + this.roomInfoList[index].n + "; 最大人数: " + this.roomInfoList[index].max;
                this.roomTextList[i].setText(strTitle);
            }
            else
            {
                this.roomList[i].visible = false;
                this.roomTextList[i].visible = false;
            }
        }
    },

    selectRoom:function()
    {
        var text = "房间名: " + this.roomID;
        game.state.states["LoginState"].leRoomName.setText(text);
        game.state.states["LoginState"].currentSelectRoomID = this.roomID;
    },

    clickPrev:function()
    {
        if(this.currentPage == 0)
        {
            return;
        }

        this.currentPage--;
        game.betApi.getRoomList();
    },

    clickNext:function()
    {
        this.currentPage++;
        game.betApi.getRoomList();
    },

    clickRename:function()
    {
        setCookie("name", "");
        location.reload();
    },

    clickLogin:function()
    {
        if(this.currentSelectRoomID < 0)
        {
            alert("请选择一个房间!");
            return;
        }

        game.betApi.setRoomID(this.currentSelectRoomID);
        game.state.start("MainState");
    },

    clickCreate:function()
    {
        game.betApi.createRoom("", 5, 10, 30, 9);
    }
};