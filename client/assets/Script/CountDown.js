var Common = require('Common');
cc.Class({
    extends: Common,
    properties: {
        //倒计时
        countdown_execution_interval:0.1,//执行间隔
        countdown_repeat_num:0,//重复了多少次
        countdown_long_time:0,//倒计时时长，单位:秒
        countdown_cycle_time:15,//倒计时一圈的时间，单位:秒
        countdown_node:{
            default:null,
            type:cc.Node
        },//倒计时，遮罩层所在的node
        countdown_task:null,//倒计时执行任务,是一个function，方便销毁定时器
        countdown_over_task:null,//倒计时结束，是一个function，执行完毕，最后会执行这个
        //动作列表
        actions:null,
        //正在进行的动作
        i:0,
        //该牌局的是数据
        hand_data:null
    },
    //添加倒计时的进度条
    add_countdown:function(seat_number,long_time){
        var me = this;
        var action_num = this.i-1;//当前第几个动作
        if(parseInt(long_time) == 0 ){
            this.scheduleOnce(function(){
                me.scheduleOnce(this.countdown_over_task,0);
            },1);
            var action_type = this.open_fixedThinkTime == 1?1:0;//动作状态 0-正常1-快进
            //判断是否有正在进行的其他动作
            this.scheduleOnce(function(){
                this.start_others(action_type,0,action_num);
            },0.5);
            return false;
        }else{
            //判断是否开启固定思考模式
            if(me.open_fixedThinkTime == 1){
                long_time = me.fixedThinkTime;
            }
        }
        var node_table_bg = cc.find("Canvas/table_bg");
        //初始化
        this.countdown_repeat_num = 0;
        this.countdown_long_time = 0;
        this.countdown_cycle_time = 15;
        if(this.countdown_task!= null){
            this.unschedule(this.countdown_task);//删除定时任务
            this.countdown_node.destroy();//删除该节点
        }
        var seat = node_table_bg.getChildByName("seat_"+seat_number);//获取该作为的节点

        //倒计时遮罩层
        var node = new cc.Node();
        node.scale = 0.85;
        //node.setPosition(seat_position);
        //node.parent = node_table_bg;
        node.parent = seat;
        node.setPosition(0.5,0.5);
        node.setLocalZOrder(3);
        var sprite = node.addComponent(cc.Sprite);
        sprite.type = cc.Sprite.Type.FILLED;
        sprite.fillType = cc.Sprite.FillType.RADIAL;
        sprite.fillCenter = new cc.Vec2(0.5,0.5);
        sprite.fillStart = 0.25;
        sprite.fillRange = 0;
        var frame = this.GameMain.getSpriteFrame("game_progress_frame");
        sprite.spriteFrame = frame;
        //倒计时文字
        var node2 = new cc.Node();
        node2.name = "time";
        var label = node2.addComponent(cc.Label);
        label.string = this.countdown_cycle_time + "s";
        label.fontSize = 30;//设置字体大小
        node2.parent = node;
        //node2.color = new cc.Color(0, 0, 0);//设置字体颜色
        node2.setPosition(0,-10);

        this.countdown_node = node;
        this.countdown_long_time = long_time;//执行时长
        this.countdown_task = function(){
            this.start_countdown(1,action_num);
        };
        this.schedule(this.countdown_task,this.countdown_execution_interval);
    },
    //倒计时开始 direction(方向):1-顺时针2-逆时针
    start_countdown:function(direction,action_num){
        this.countdown_repeat_num++;//已经执行了多少次
        var full_cycle = 1;//一个整圈
        var speed = full_cycle / this.countdown_cycle_time;//速度

        var sprite = this.countdown_node.getComponent(cc.Sprite);
        var fillRange = sprite.fillRange;
        var cost_time = this.countdown_execution_interval * this.countdown_repeat_num;//耗时
        if(parseInt(cost_time) == cost_time ){
            var time_node = this.countdown_node.getChildByName("time");
            var time_label = time_node.getComponent(cc.Label);
            time_label.string =  (parseInt(time_label.string) - 1)+"s";
        }

        if(direction == 2){
            //逆时针
            fillRange = cost_time < this.countdown_long_time ? fillRange += (this.countdown_execution_interval * speed):0;
        }else{
            //顺时针
            fillRange = cost_time < this.countdown_long_time ? fillRange -= (this.countdown_execution_interval * speed):0;
        }
        sprite.fillRange = fillRange;
        //如果开启固定思考时间的设置，并且这个过程的执行时间大于思考时间时，直接结束倒计时
        var action_type = 0;//动作状态 0-正常1-快进
        if(this.open_fixedThinkTime == 1 && this.countdown_long_time >= this.fixedThinkTime){
            fillRange = 0;
            action_type = 1;
        }
        //判断是否有正在进行的其他动作
        this.start_others(action_type,cost_time,action_num);

        if(fillRange == 0){
            this.unschedule(this.countdown_task);//删除定时任务
            this.countdown_node.destroy();//删除该节点
            this.countdown_task = null;
            //如果执行完毕，判断是否有下一步动作，如果有，执行下一步
            if(this.countdown_over_task != null){
                this.scheduleOnce(this.countdown_over_task,0);
            }
        }
    },
    //延长时间
    delay_countdown:function(seat_number,duration){
        var remain_time = parseInt(this.countdown_cycle_time - duration);//剩余的时间
        var total_time = 15 + remain_time;
        this.countdown_cycle_time = total_time; //一圈的总耗时
        this.countdown_repeat_num = 0; //已经执行了多少
        if(this.countdown_task!= null){
            this.unschedule(this.countdown_task);//删除定时任务
            this.countdown_node.destroy();//删除该节点
        }
        //延时
        var node_table_bg = cc.find("Canvas/table_bg");
        var seat = node_table_bg.getChildByName("seat_"+seat_number);//获取该作为的节点
        //倒计时遮罩层
        var node = new cc.Node();
        node.scale = 0.85;
        node.parent = seat;
        node.setPosition(0.5,0.5);
        node.setLocalZOrder(3);
        var sprite = node.addComponent(cc.Sprite);
        sprite.type = cc.Sprite.Type.FILLED;
        sprite.fillType = cc.Sprite.FillType.RADIAL;
        sprite.fillCenter = new cc.Vec2(0.5,0.5);
        sprite.fillStart = 0.25;
        sprite.fillRange = 0;
        var frame = this.GameMain.getSpriteFrame("game_progress_frame");
        sprite.spriteFrame = frame;
        //倒计时文字
        var node2 = new cc.Node();
        node2.name = "time";
        var label = node2.addComponent(cc.Label);
        label.string = this.countdown_cycle_time + "s";
        label.fontSize = 30;//设置字体大小
        node2.parent = node;
        node2.setPosition(0,-10);

        this.countdown_node = node;
        this.countdown_task = function(){
            this.start_countdown(1);
        };
        this.schedule(this.countdown_task,this.countdown_execution_interval);
        var me = this;
        this.countdown_over_task = function(){
            me.actionend();
        };
    },
    //开始执行其他动作:type:0-正常1-快进 cost：已经执行了多少秒,action_num是之前的this.i-1
    start_others:function(type,cost_time,action_num){
        //判断等待动作时，是否有其他人进行的动作
        var current_action = this.actions[action_num];
        if(current_action!=undefined && current_action['others']!=undefined && current_action['others']!=null){
            for(var key in current_action['others']){
                var other_key = current_action['others'][key];
                var duration = current_action['duration']-(current_action['timestamp']-this.hand_data['others'][other_key]['timestamp']);
                if(duration<0){
                    cc.log("时间间隔为负，请检查数据");
                }
                if(type == 0 && parseInt(cost_time) == cost_time && cost_time == duration){
                    //正常流程，待到整点时，进行判断是否有其他的动作
                    this.other_quit(this.hand_data['others'][other_key]['chair_id']);
                }
                //duration<cost_time表示如果已经执行过了，不再执行
                if(type == 1 && duration>=cost_time){
                    //快进时，判断等待动作时，是否有其他人进行的动作(未执行过的)
                    this.other_quit(this.hand_data['others'][other_key]['chair_id']);
                }
            }
        }
    },
    //ajax 请求
    reqstart:function(){
        var hand_id = this.getQueryString("hand_id");
        var url="";
        //hand_id = 32634;
        //hand_id = 32928;
        //hand_id = 33305;
        //hand_id = 33308;
        //hand_id=33310;
        //hand_id=33509;
        if(hand_id != null){
            var host_name = window.location.host;
            var reg = /^localhost:/;
            if(reg.test(host_name)==true){
                url = "http://qa-api.kkpoker.com:8090/Html/get_mongo_data/hand_id/"+hand_id;
            }else{
                url = "http://" + host_name + "/Html/get_mongo_data/hand_id/"+hand_id;
            }
        }else{
            //测试数据
            url="http://172.16.0.210:2016/info.php";
         }


        var xhr = new XMLHttpRequest();
        xhr.open("GET", url, true);
        xhr.send();
        var me=this;
        xhr.onreadystatechange = function () {
            if (xhr.readyState == 4) {
                if (xhr.status == 200) {
                    var response = eval('(' + xhr.responseText + ')');
                    response = me.sorting_data(response);
                    var len = response["actions"].length;
                    for(var i=0;i<len;i++){
                        if(i==0){
                            response["actions"][i]["duration"]=response["actions"][i]["timestamp"]-response["start"]["timestamp"];
                        }else{
                            if(response["actions"][i-1]['CMD'] == 19){
                                response["actions"][i]["duration"] = 0;
                            }else{
                                response["actions"][i]["duration"]=response["actions"][i]["timestamp"]-response["actions"][i-1]["timestamp"];
                            }
                        }
                    }
                    //初始化动作属性
                    me.actions = response["actions"];
                    me.i=0;
                    me.hand_data = response;
                    cc.log(JSON.stringify(response));
                    me.sit_down();//坐下
                } else {

                }
            }
        }
    },

    //组装action和other
    sorting_data:function(hand_data){
        //过滤掉无用的others动作
        if(hand_data['others']!=undefined){
            var others_data=[];
            for(var key in hand_data['others']){
                if(hand_data['others'][key]['CMD'] == 5 ){
                    hand_data['others'][key]['type']='others';
                    others_data.push(hand_data['others'][key]);
                }
            }
            hand_data['others']=others_data;
        }
        var action_data=[];
        if(hand_data['actions']!=undefined){
            action_data = action_data.concat(hand_data['actions']);
        }
        if(hand_data['others']!=undefined){
            action_data = action_data.concat(hand_data['others']);
        }
        var action_data_len = action_data.length;
        //按时间排序
        for(var i=0;i<action_data_len-1;i++){
           for(var j=i+1;j<action_data_len;j++){
               if(parseInt(action_data[j]['timestamp'])<parseInt(action_data[i]['timestamp'])){
                    var a = action_data[i];
                    action_data[i]=action_data[j];
                    action_data[j]=a;
               }
           }
        }
        //找到同时进行的工作
        var action_result=[];//动作列表
        var other_result=[];//同时进行的动作
        for(var i=0;i<action_data_len;i++){
            var other_is_sum_up=0;//是否需要归纳0-否1-是
            var data = null;//命令数据
            if(action_data[i]['type'] == "others"){
                data={"CMD":9999,"chair_id":action_data[i]['current_action_chair'],"chip":0,"pot":0,"timestamp":action_data[i]['timestamp']}
                other_result.push(action_data[i]);//把同时进行的动作存起来
                other_is_sum_up=1;
            }else{
                data = action_data[i];
            }
            var action_result_len = action_result.length;
            if(action_result_len>0 && data['CMD']==9999 && data['chair_id'] == action_result[action_result_len-1]['chair_id']){
                action_result[action_result_len-1]['timestamp']=data['timestamp'];
                data = null;
            }
            //有可能同时做的动作，其他人都没正在做动作
            if(data != null && data['chair_id']!=0){
                action_result.push(data);
            }
            if(other_is_sum_up == 1){
                var sort_num = action_result.length - 1;
                //把other动作绑定到action动作当中
                if(action_result[sort_num]['others']==undefined){
                    action_result[sort_num]['others']=[other_result.length-1];
                }else{
                    action_result[sort_num]['others'].push(other_result.length-1);
                }
            }
        }
        //再次归纳action动作
        var action_result2 = [];
        for(var i=0;i<action_result.length;i++){
            var is_in=1;//是否放入action_result2
            if(action_result[i+1]!=undefined && action_result[i+1]!=null){
                //如果该动作与下一个动作都是9999，且chair_id相等的话，归纳到一起，或者下一个正在进行的是99的话，修改下一个的chair_id为上一个
                if(action_result[i]['CMD']==9999 && (action_result[i]['chair_id'] == action_result[i+1]['chair_id'] ||  action_result[i+1]['chair_id']==99)){
                    action_result[i+1]['chair_id']=action_result[i+1]['chair_id']==99?action_result[i]['chair_id']:action_result[i+1]['chair_id'];
                    action_result[i+1]['others']=action_result[i+1]['others']?action_result[i+1]['others']:[];
                    action_result[i+1]['others']=action_result[i]['others'].concat(action_result[i+1]['others']);
                    is_in=0;
                }
            }else if(action_result[i]['CMD']==9999 && action_result[i]['chair_id']==99 & i>0){
                if(action_result2.length>0){
                    action_result[i]['chair_id'] = action_result2[action_result2.length-1]['current_action_chair'];
                    //修正：如果发现chair_id=99，无效的操作人，找到这个命令的其他操作的人，座位为当前这个人
                    if(action_result[i]['chair_id']==99 && action_result[i]['others']!=undefined){
                        action_result[i]['chair_id'] = hand_data['others'][action_result[i]['others'][0]]['chair_id'];
                    }
                }
            }else if(action_result[i]['CMD']==9999 && action_result[i]['chair_id']==99 & action_result.length == 1){
                //修正：当只有一个9999的动作时，当时正在进行的动作是小盲在思考
                action_result[i]['chair_id'] = hand_data['start']['sb_chair'];
            }

            if(is_in==1){
                action_result2.push(action_result[i]);
            }
        }
        hand_data['actions']=action_result2;
        hand_data['others']=other_result;
        return hand_data;
    }
});
