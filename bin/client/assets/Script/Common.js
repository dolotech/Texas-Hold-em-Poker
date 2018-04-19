cc.Class({
    extends:  cc.Component,

    properties: {
        seat:{
            default:[],
            type:cc.Node
        },
        Lang:null,//语言版本
        SourceSuffix:null,//资源后缀
        GameMain:null,//游戏资源
        GameCards:null,//游戏的牌中的资源
        audio_chipsToTable:null,//下注的音频
        audio_check:null,//check的音频
        audio_fold:null,//fold的音频
        audio_distributeCard:null,//翻牌的声音
        open_fixedThinkTime:0,//开启固定思考时间的设置0-否1-开启
        fixedThinkTime:2,//固定思考时间，配合open_fixedThinkTime使用，单位：秒
        open_mute:0,//是否开启静音模式0-否1-是
        is_mobile:0,//是否是手机访问0-否1-是
        fontStyle:null,//字体样式
        //websocket
        ws:null,
    },
    //web socket连接
    wsstart:function(){
        var me=this;
        var ws = new WebSocket("ws://127.0.0.1:3564");
        ws.binaryType = "arraybuffer" ;
        var decoder = new TextDecoder('utf-8')


        ws.onopen = function (event) {
            me.ws = ws;
            console.log("Send Text WS was opened.");
        };

        ws.onmessage = function (event) {
            console.log("response text msg: " + event.data);

            var input = event.data;

            var data = JSON.parse(decoder.decode(event.data));
            console.log("server response");

            console.log(data);

        };
        ws.onerror = function (event) {
            console.log("Send Text fired an error");
        };
        ws.onclose = function (event) {
            console.log("WebSocket instance closed.");
        };

        //setTimeout(function () {
        //    if(me.ws!=null){
        //        cc.log(me.ws.readyState);
        //        if (me.ws.readyState === WebSocket.OPEN) {
        //            cc.log("send message to server3000");
        //            me.ws.send(JSON.stringify({Ai: {
        //                Ai_id: '123456',
        //                Hand:[1,2]
        //            }}))
        //        }
        //        else {
        //            console.log("WebSocket instance wasn't ready...");
        //        }
        //    }else{
        //        cc.log("ws is null");
        //    }
        //}, 3000);

    },
    sit_down:function(){
        //加载游戏资源图片
        var me = this;
        me.is_mobile = this.isMobile();
        if(me.is_mobile == 1){
            me.fontStyle={
                "table_name":{"fontSize":30,"lineHeight":30},
                "table_sb":{"fontSize":30,"lineHeight":30},
                "table_code":{"fontSize":30,"lineHeight":30},
                "nick":{"fontSize":30,"lineHeight":30},
                "chips":{"fontSize":30,"lineHeight":30},
                "chip":{"fontSize":30,"lineHeight":30},
                "pot":{"fontSize":30,"lineHeight":30}
            };
        }else{
            me.fontStyle={
                "table_name":{"fontSize":25,"lineHeight":25},
                "table_sb":{"fontSize":25,"lineHeight":25},
                "table_code":{"fontSize":25,"lineHeight":25},
                "nick":{"fontSize":25,"lineHeight":25},
                "chips":{"fontSize":25,"lineHeight":25},
                "chip":{"fontSize":25,"lineHeight":25},
                "pot":{"fontSize":25,"lineHeight":25}
            };
        }
        cc.loader.loadRes("game_cards",cc.SpriteAtlas,function(err,atlas){
            me.GameCards = atlas;
        });
        //加载游戏音频
        if(this.open_mute == 0){
            cc.loader.loadRes("audio/audio_chipsToTable", function (err, assets) {
                me.audio_chipsToTable = assets;
            });
            cc.loader.loadRes("audio/audio_check", function (err, assets) {
                me.audio_check = assets;
            });
            cc.loader.loadRes("audio/audio_fold", function (err, assets) {
                me.audio_fold = assets;
            });
            cc.loader.loadRes("audio/audio_distributeCard", function (err, assets) {
                me.audio_distributeCard = assets;
            });
        }

        var table_data = this.hand_data;
        table_data['table_name']=table_data['table_name']?table_data['table_name']:"";
        table_data['table_code']=table_data['table_code']?table_data['table_code']:"";
        var node_table_bg = cc.find("Canvas/table_bg");

        //异步加载头像，不能放在循环内
        var load_avatar = function(url,sprite_user){
            cc.loader.load(url,function(err,tex){
                var frame  = new cc.SpriteFrame(tex,cc.Rect(0, 0, 87, 123));
                sprite_user.spriteFrame = frame;
            });
        };
        cc.loader.loadRes("GameMain_6p",cc.SpriteAtlas,function(err,atlas){
            me.GameMain = atlas;
            var sb_data = null;//小盲的数据
            var bb_data = null;//大盲的数据
            //坐下
            for(var k in table_data['players']){
                var v = table_data['players'][k];
                var seat_node = node_table_bg.getChildByName("seat_"+v['chair_id']);
                var node_size = seat_node.getContentSize();//获取node的尺寸
                //名字
                var label_nick = seat_node.getChildByName("nick").getComponent(cc.Label);
                label_nick.string = me.getLength(v['nick'])>10?me.cutStr(v['nick'],7):v['nick'];
                label_nick.fontSize = me.fontStyle['nick']['fontSize'];
                label_nick.lineHeight = me.fontStyle['nick']['lineHeight'];

                //带入的筹码
                var label_chips = seat_node.getChildByName("chips").getComponent(cc.Label);
                label_chips.string = v['remain_chip']+v['table_chip'];
                label_chips.fontSize = me.fontStyle['chips']['fontSize'];
                label_chips.lineHeight = me.fontStyle['chips']['lineHeight'];
                //头像
                var node_mark = new cc.Node();
                node_mark.name='avatar';
                var mask_user = node_mark.addComponent(cc.Mask);
                mask_user.type = cc.Mask.ELLIPSE;
                node_mark.width = node_size['width']-15;
                node_mark.height = node_size['height']-15;
                node_mark.setPosition(0.5,0.5);
                node_mark.parent =  seat_node;
                
                var node_user = new cc.Node();
                var sprite_user = node_user.addComponent(cc.Sprite);
                node_user.parent = node_mark;
                if(v['avatar'] != null && v['avatar']!=""){
                    node_user.scale = (node_size['width']-15)/120;
                    load_avatar(v['avatar'],sprite_user);
                }else{
                    sprite_user.spriteFrame =  me.GameMain.getSpriteFrame("game_seat_valid");
                }

                seat_node.getChildByName("game_tip").setLocalZOrder(2);

                if(v['chair_id'] == table_data['start']['sb_chair']){
                    sb_data = v;
                }else if(v['chair_id'] == table_data['start']['bb_chair']){
                    bb_data = v;
                }
            }
            //牌局名字
            var label = node_table_bg.getChildByName("table_name").getComponent(cc.Label);
            label.string = "§ " + table_data['table_name'] + " §";
            label.fontSize = me.fontStyle['table_name']['fontSize'];
            label.lineHeight = me.fontStyle['table_name']['lineHeight'];

            //牌局号,如果是快速牌局显示牌局号，如果是俱乐部牌局，显示盲注
            if(parseInt(table_data['table_code'])!=0){
                var label_table_code = node_table_bg.getChildByName("table_code").getComponent(cc.Label);
                label_table_code.string = table_data['table_code'];
                label_table_code.fontSize = me.fontStyle['table_code']['fontSize'];
                label_table_code.lineHeight = me.fontStyle['table_code']['lineHeight'];
            }
            //盲注
            var label_table_sb = node_table_bg.getChildByName("table_sb").getComponent(cc.Label);
            label_table_sb.string = me.ConvertLang("blind")+" "+sb_data['table_chip']+"/"+bb_data['table_chip'];
            label_table_sb.fontSize = me.fontStyle['table_sb']['fontSize'];
            label_table_sb.lineHeight = me.fontStyle['table_sb']['lineHeight'];

            //删除加载图片
            var splash = document.getElementById('splash');
            if(undefined!=splash){
                splash.style.display = 'none';
            }
        });

    },
    //监听座位点击事件
    onSeat:function(){
        cc.log("onSeat init");
        var node_table_bg = cc.find("Canvas/table_bg");
        for(var i=0;i<9;i++){
            //var v = table_data['players'][k];
            var seat_node = node_table_bg.getChildByName("seat_"+i);
            seat_node.on("mouseup",function(event){
                this.seatClick (event);
            },this);
        }
    },
    //点击座位坐下 未完成
    seatClick:function(e){
        var me = this;

        var rName = cc.random0To1();
        var v={};
        v["nick"]   =   "波霸"+rName;
        v["remain_chip"]    =   "1000";
        v["table_chip"]     =   "10";
        //v["avatar"]         =   "/1_5734451068aaa.jpg";
        v["avatar"]         =   "";


        if(me.ws){
            //请求后台坐下
            if (me.ws.readyState === WebSocket.OPEN) {
                cc.log("send message to server3000");
                me.ws.send(JSON.stringify({C2S_Join: {
                                    UserId: '123456',
                                    UserName:v["nick"],
                                    TableId:'22222',
                                    UserToken:'33333',
                                }}));
                me.ws.send(JSON.stringify({C2S_Login: {
                    UserAccount: "tian",
                    UserName: "nick",
                    UserPassword: "123456",
                }}));
            }
            else {
                this.wsstart();
                console.log("WebSocket instance wasn't ready...");
            }
        };
        var seatName = e.currentTarget._name;


        var node_table_bg = cc.find("Canvas/table_bg");
        var seat_node = node_table_bg.getChildByName(seatName);
        var node_size = seat_node.getContentSize();//获取node的尺寸
        //名字
        var label_nick = seat_node.getChildByName("nick").getComponent(cc.Label);
        label_nick.string = me.getLength(v['nick'])>10?me.cutStr(v['nick'],7):v['nick'];
        label_nick.fontSize = 30;
        label_nick.lineHeight = 30;

        //带入的筹码
        var label_chips = seat_node.getChildByName("chips").getComponent(cc.Label);
        label_chips.string = v['remain_chip']+v['table_chip'];
        label_chips.fontSize = 30;
        label_chips.lineHeight = 30;

        //头像
        var node_mark = new cc.Node();
        node_mark.name='avatar';
        var mask_user = node_mark.addComponent(cc.Mask);
        mask_user.type = cc.Mask.ELLIPSE;
        node_mark.width = node_size['width']-15;
        node_mark.height = node_size['height']-15;
        node_mark.setPosition(0.5,0.5);
        node_mark.parent =  seat_node;

        var node_user = new cc.Node();
        var sprite_user = node_user.addComponent(cc.Sprite);
        node_user.parent = node_mark;
        if(v['avatar'] != null && v['avatar']!=""){
            //异步加载头像，不能放在循环内
            var load_avatar = function(url,sprite_user){
                cc.loader.load(url,function(err,tex){
                    var frame  = new cc.SpriteFrame(tex,cc.Rect(0, 0, 87, 123));
                    sprite_user.spriteFrame = frame;
                });
            };
            node_user.scale = (node_size['width']-15)/120;
            load_avatar(v['avatar'],sprite_user);
        }else{
            if(this.GameMain == null){
                cc.loader.loadRes("GameMain_6p",cc.SpriteAtlas,function(err,atlas){
                    this.GameMain = atlas;
                    sprite_user.spriteFrame =  this.GameMain.getSpriteFrame("game_seat_valid");
                });
            }else{
                sprite_user.spriteFrame =  this.GameMain.getSpriteFrame("game_seat_valid");
            }
        }
        seat_node.getChildByName("game_tip").setLocalZOrder(2);
    },
    //js获取当前url的参数
    getQueryString:function(name){
        if(undefined==window.location){
            return null;
        };
        var reg = new RegExp("(^|&)"+ name +"=([^&]*)(&|$)");
        var r = window.location.search.substr(1).match(reg);
        if(r!=null){
            return  unescape(r[2]);
        }else{
            return null;
        }
    },
    //获取字符串长度
    getLength:function (str) {
        ///<summary>获得字符串实际长度，中文2，英文1</summary>
        ///<param name="str">要获得长度的字符串</param>
        var realLength = 0, len = str.length, charCode = -1;
        for (var i = 0; i < len; i++) {
            charCode = str.charCodeAt(i);
            if (charCode >= 0 && charCode <= 128) realLength += 1;
            else realLength += 2;
        }
        return realLength;
    },
    //s截取字符串，中英文都能用@param str：需要截取的字符串@param len: 需要截取的长度
    cutStr:function (str, len) {
        var str_length = 0;
        var str_len = str.length;
        var str_cut = new String();
        for (var i = 0; i < str_len; i++) {
            var a = str.charAt(i);
            str_length++;
            if (escape(a).length > 4) {
                //中文字符的长度经编码之后大于4
                str_length++;
            }
            str_cut = str_cut.concat(a);
            if (str_length >= len) {
                str_cut = str_cut.concat("...");
                return str_cut;
            }
        }
        //如果给定字符串小于指定长度，则返回源字符串；
        if (str_length < len) {
            return str;
        }
    },

    //开启静音模式status:0-关闭1-开启
    set_mute:function(){
        var sprite_url = "";
        if(this.open_mute == 1){
            this.open_mute = 0;
            sprite_url='open_audio';
        }else{
            this.open_mute = 1;
            sprite_url='close_audio';
        }
        var fast_forward =cc.find("Canvas/sound");
        var sprite = fast_forward.getComponent(cc.Sprite);
        sprite.spriteFrame = this.GameMain.getSpriteFrame(sprite_url);
    },
    //判断是否是手机访问
    isMobile:function(){
        var ua = navigator.userAgent;
        var ipad = ua.match(/(iPad).*OS\s([\d_]+)/),
        isIphone = !ipad && ua.match(/(iPhone\sOS)\s([\d_]+)/),
        isAndroid = ua.match(/(Android)\s+([\d.]+)/),
        isMobile = isIphone || isAndroid;
        var bool = isMobile === null?0:1;
        return bool;
    },
    //语言翻译
    ConvertLang:function(key){
        var config_lang={};
        config_lang['zh-cn']={"pot":"底池","blind":"盲注","mute":"静音","no_think":"忽略思考时间"};
        config_lang['zh-tw']={"pot":"底池","blind":"盲注","mute":"靜音","no_think":"忽略思考時間"};
        config_lang['thai-th']={"pot":"กองกลาง","blind":"Stakes","mute":"ปิดเสียง","no_think":"fast-forwward"};
        config_lang['en-us']={"pot":"Pot","blind":"Stakes","mute":"Mute","no_think":"fast-forwward"};
        config_lang['ko']={"pot":"바닥풀 ","blind":"맹주","mute":"정음 ","no_think":"忽略思考时间"};
        return config_lang[this.Lang][key];
    },
    //多语言替换座位
    //ReplaceSeat:function(){
    //    if(this.Lang == 'thai-th'){
    //        for(var i=0;i<9;i++){
    //            cc.find("Canvas/table_bg/seat_"+i).getComponent(cc.Sprite).setVisible(false);
    //        }
    //        cc.loader.loadRes("GameMain_th_6p",cc.SpriteAtlas,function(err,atlas){
    //            var game_seat_empty = atlas.getSpriteFrame("game_seat_empty");
    //            for(var i=0;i<9;i++){
    //                var sprite = cc.find("Canvas/table_bg/seat_"+i).getComponent(cc.Sprite)
    //                sprite.spriteFrame = game_seat_empty;
    //                sprite.setVisible(true);
    //            }
    //        });
    //    }
    //}
});