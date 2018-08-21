define(["jquery", "amrnb", "ionRangeSlider"], function($) {
	"use strict";
	var hasPlay = false; //是否播放过的标识

	var Player = function(){

		var volumeFlag = false;  //是否静音的标识

		this.baseId = null;
		this.playUrl = null;
		this.pauseUrl = null;
		this.stopUrl = null;
		this.audioUrl = null;
		this.progressColor = null;
		this.playBarColor = null;
		this.audioType = null;
		this.supportType = "mp3,wav";

		/**
		 * 初始化录音组件
		 * @param playUrl   播放按钮图片路径
		 * @param pauseUrl  暂停按钮图片路径
		 * @param stopUrl   停止按钮图片路径
		 * @param baseId    容器id
		 * @param audioUrl  录音文件路径
		 * @param progressColor   设置播放进度条颜色
		 * @param playBarColor   设置播放条背景色
		 */
		this.config = function(playUrl, pauseUrl, stopUrl, baseId, audioUrl, progressColor, playBarColor){
			this.baseId = baseId;
			this.playUrl = playUrl;
			this.pauseUrl = pauseUrl;
			this.stopUrl = stopUrl;
			this.audioUrl = audioUrl;
			this.progressColor = progressColor;
			this.playBarColor = playBarColor;
			var pos= audioUrl.lastIndexOf(".");//查找最后一个.的位置,截取文件扩展名
			this.audioType = $.trim(audioUrl.substring(pos + 1).toLowerCase());
			this.genPlayer(playUrl, stopUrl);
		};

		this.genPlayer = function(playUrl, stopUrl){
			var playerContent = '<div id="m_playBar" class="sample-amr box-shadow">'+
									'<div id="m_playImgs">'+
										'<img id="m_play" class="play" src="'+ playUrl+ '" style="margin-right: 15px;">'+
										'<img id="m_stop" class="stop" src="'+ stopUrl+ '">'+
									'</div>'+
									'<span class="curTime">00:00:00</span>'+
									'<div id="m_playSlider" style="flex-shrink: 1;flex-grow: 2;">'+
										'<input  class="playSlider" type="hidden" name="playSlider" value="" />'+
									'</div>'+
									'<span class="totalTime">00:00:00</span>'+
									'<div id="m_playSound">'+
										'<img id="m_soundImg" src="' + MyPlayer.getBaseUrl() + '/img/sound.png">'+
									'</div>'+
									'<div id="m_playVolume" class="v">'+
										'<input  class="volumeSlider" type="hidden" name="volumeSlider" value="" />'+
									'</div>'+
								'</div>';

			$("#" + this.baseId).html(playerContent);

			//设置录音播放条的高度
			var playBarHeight = $("#" + this.baseId).height() > 39 ? $("#" + this.baseId).height() : 45;

			//设置音量调的上下padding值，使其居中
			var playVolumePaddingTB = parseInt((playBarHeight - 12) / 2);
			$("#m_playVolume").css({paddingTop: playVolumePaddingTB + 'px', paddingBottom: playVolumePaddingTB + 'px'});

			//设置播放条样式
			$("#m_playBar").css({background: this.playBarColor, lineHeight: playBarHeight + 'px'});

			//设置按钮行高，使其居中
			$("#m_playImgs").css({lineHeight: (playBarHeight - 3) + 'px'});

			//初始化播放条和音量条
			this.initEventHandler();

			//初始化隐藏播放进度条按钮
			$("#m_irs-slider").hide();

			//设置播放进度条的颜色
			$(".irs-bar,.irs-bar-edge").css({
				borderColor: this.progressColor,
				background: this.progressColor
			});
		};

		/**
		 * 初始化播放条和音量条
		 */
		this.initEventHandler = function(){
			var mp = this;

			//播放条
			$("#" + mp.baseId + " .playSlider").ionRangeSlider({
				type: 'single',
				min: 0,
				hide_min_max: true,
				hide_from_to: true,
				onChange:function(obj){

					//播放之前播放进度条点击不能移动
					if(!hasPlay) {
						$("#m_playSlider #m_irs-bar").width(0);
						return;
					}

					//设置播放进度条的颜色
					$(".irs-bar,.irs-bar-edge").css({
						borderColor: mp.progressColor,
						background: mp.progressColor
					});

					//播放时拖动进度条
					if(boxes.status=="play"){
						boxes.stop();
					}
					console.log("当前值:"+obj.from);
					boxes.$curTime.text(boxes.formatTime(obj.from));
					boxes.$totalTime.text(boxes.formatTime(obj.max));
				}
			});

			var boxes = new Boxes();
			boxes.playUrl = mp.playUrl;
			boxes.pauseUrl = mp.pauseUrl;
			boxes.baseId = mp.baseId;
			boxes.progressColor = mp.progressColor;

			/**
			 * 切换音量按钮
			 */
			$("#m_soundImg").click(function(){
				var slider = $("#" + boxes.baseId + " .volumeSlider").data("ionRangeSlider");
				if(!volumeFlag) {
					boxes.volumeOldValue = boxes.volumeValue;
					$("#m_soundImg").attr("src", MyPlayer.getBaseUrl() + "/img/sound2.png");
					slider.update({
						from:0
					});
					volumeFlag = true;
				} else {
					boxes.volumeValue = boxes.volumeOldValue;
					$("#m_soundImg").attr("src", MyPlayer.getBaseUrl() + "/img/sound.png");
					slider.update({
						from: boxes.volumeValue
					});
					volumeFlag = false;
					$(".irs-bar,.irs-bar-edge").css({
						borderColor: boxes.progressColor,
						background: boxes.progressColor
					});
				}

				boxes.volumeValue = parseFloat($("#"+boxes.baseId+" .volumeSlider").val());
				if(boxes.gainNode){
					boxes.gainNode.gain.value = boxes.volumeValue;
				}
			});

			//音量条
			$("#" + this.baseId + " .volumeSlider").ionRangeSlider({
				type: 'single',
				min: 0,
				max: 1,
				step: 0.1,
				from: boxes.volumeValue,
				hide_min_max: true,
				hide_from_to: true,
				force_edges: true,
				onChange: function() {
					$(".irs-bar,.irs-bar-edge").css({
						borderColor: boxes.progressColor,
						background: boxes.progressColor
					});
				},
				onFinish: function(obj){
					var volumeValue = $("#" + boxes.baseId + " .volumeSlider").val();
					if(volumeValue == "0") {
						$("#m_soundImg").attr("src", MyPlayer.getBaseUrl() + "/img/sound2.png");
						volumeFlag = true;
					} else {
						volumeFlag = false;
						$("#m_soundImg").attr("src", MyPlayer.getBaseUrl() + "/img/sound.png");
					}

					boxes.volumeValue = parseFloat($("#" + boxes.baseId + " .volumeSlider").val());
					if(boxes.gainNode){
						boxes.gainNode.gain.value = boxes.volumeValue;
					}
				}
			});

			boxes.initAPI();
			boxes.$curTime = $("#" + boxes.baseId + " .sample-amr .curTime");
			boxes.$totalTime = $("#" + boxes.baseId + " .sample-amr .totalTime");

			//播放暂停
			E("#" + mp.baseId + " .sample-amr .play").onclick = function() {
				//保存按钮引用
				boxes.playBtn = this;
				if(boxes.status == "stop"){
                    if(mp.supportType.indexOf(mp.audioType)==-1){
                        alert("不支持的录音类型");
                        return;
                    }
                    this.src = boxes.pauseUrl;
					boxes.fetchBlob(mp.audioUrl, function(blob) {
						//这之内的数据解码和加载都是需要时间的,故播放操作要写在其内的回调中
						if(mp.audioType=="mp3"){
							boxes.playAmrArray(blob,"mp3");
						}else if(mp.audioType=="wav"){
							boxes.playAmrBlob(blob);
						}else{
							alert("类型异常");
							return;
						}
					},mp.audioType);
				}else if(boxes.status == "play"){
					this.src = boxes.playUrl;
					boxes.playStart = parseFloat($("#"+boxes.baseId+" .playSlider").val());
					boxes.stop();
				}else{
					alert("状态异常");
				}
			};

			//终止
			E("#" + mp.baseId + " .sample-amr .stop").onclick = function() {
                if(mp.supportType.indexOf(mp.audioType)==-1){
                    alert("不支持的录音类型");
                    return;
                }
				if(boxes.playBtn && boxes.playBtn.src) {
					boxes.playBtn.src = boxes.playUrl;
				} else {
					return;
				}
				boxes.terminal();
			}
		};
	}

	function E(selector) {
		return document.querySelector(selector);
	}

	var Boxes = function() {
		this.baseId = null;
		this.playUrl = null;
		this.pauseUrl = null;
		this.playBtn = null;     //播放按钮引用
		this.interId = null;     //定时调度id
		this.status = "stop"; //当前播放状态
		this.startTime = 0;  //开始播放时间
		this.totalTime = 0;  //总时间
		this.totalSeconds = 0; //总秒数
		this.$curTime = null;  //显示当前时间的jquery对象
		this.$totalTime = null;  //显示总时间的jquery对象
		this.playStart = 0;  //音频文件起始播放时间点
		this.buffer = null; //要处理的文件
		this.audioContext = null; //进行音频处理的上下文，稍后会进行初始化
		this.source = null; //保存音频
		this.gainNode = null; //音频处理节点
		this.volumeValue = 0.5; //音量值
		this.volumeOldValue = null; //静音之前的音量值

		this.initAPI= function() {
			//统一前缀，方便调用
			window.AudioContext = window.AudioContext || window.webkitAudioContext || window.mozAudioContext || window.msAudioContext;
			//安全地实例化一个AudioContext并赋值到audioContext属性上，方便后面处理音频使用
			try {
				this.audioContext = new AudioContext();
			} catch (e) {
				alert("您的浏览器不支持录音组件，建议使用(谷歌/火狐)浏览器访问");
				console.log(e);
			}
		};

		this.fetchBlob= function(url, callback,audioType) {
			
			var xhr = new XMLHttpRequest();
			xhr.open('GET', url,true);
			if(audioType=="mp3"){
				xhr.responseType = 'arraybuffer';
			}else if(audioType=="wav"){
				xhr.responseType = 'blob';
			}else{
				alert("unknown audio type");
				return;
			}
			xhr.onload = function() {
				callback(this.response);
			};
			xhr.onerror = function() {
				alert('Failed to fetch ' + url);
			};
			xhr.send();
			
		};

		this.readBlob= function(blob, callback) {
			var reader = new FileReader();
			reader.onload = function(e) {
				var data = new Uint8Array(e.target.result);
				callback(data);
			};
			reader.readAsArrayBuffer(blob);
		};

		this.playAmrBlob= function(blob, callback) {
			var self = this;
			this.readBlob(blob, function(data) {
				self.playAmrArray(data,"wav");
			});
		};

		this.playAmrArray= function(array,audioType) {
			var self = this;
			if(audioType=="mp3"){
				this.audioContext.decodeAudioData(array).then(function(decodedData) {
					// use the decoded data here
					if (!decodedData) {
						alert('Failed to decode!');
						return;
					}
					self.playPcm(decodedData,"mp3");
				});
			}else if(audioType=="wav"){
				var samples = AMR.decode(array);
				if (!samples) {
					alert('Failed to decode!');
					return;
				}
				this.playPcm(samples,"wav");
			}else{
				alert("unknown audio type");
				return;
			}

		};

		this.playPcm= function(samples,audioType) {
			this.source = this.audioContext.createBufferSource();  //创建一个声音源
			this.gainNode = this.audioContext.createGain();  //创建一个gain node
			if(audioType=="mp3"){
				this.buffer = samples;
			}else if(audioType=="wav"){
				if(this.buffer){

				}else{
					this.buffer = this.audioContext.createBuffer(1, samples.length, 8000);
					if (this.buffer.copyToChannel) {
						this.buffer.copyToChannel(samples, 0, 0)
					} else {
						var channelBuffer = this.buffer.getChannelData(0);
						channelBuffer.set(samples);
					}
				}
			}else{
				alert("unknown audio type");
				return;
			}

			this.source.buffer = this.buffer;
			this.source.connect(this.gainNode);  //将实例与gain node相连
			this.gainNode.connect(this.audioContext.destination);  //将gain node与播放设备连接
				this.gainNode.gain.value = this.volumeValue;   //控制音量
				this.totalTime = this.formatTime(this.buffer.duration);
				this.totalSeconds = this.buffer.duration;
				console.log("totalTime:" + this.totalTime)
				this.$totalTime.text(this.totalTime);
				this.play();
			};

		//停止(暂停)
		this.stop = function(){
			if(this.status=="play"){
				this.stopPro();
				this.source.stop();
				//重新初始化一个source
				this.source = this.audioContext.createBufferSource();
				this.gainNode = this.audioContext.createGain();
				this.source.buffer = this.buffer;
				this.source.connect(this.gainNode);
				this.gainNode.connect(this.audioContext.destination);
				this.gainNode.gain.value = this.volumeValue;

				this.status="stop";
			}
		};

		//终止
		this.terminal=function(){
			this.stop();
			this.playStart = 0;
			var slider = $("#"+this.baseId+" .playSlider").data("ionRangeSlider");
			slider.update({
				from:0
			});
			this.$curTime.text("00:00:00");
			this.$totalTime.text("00:00:00");
		};

		//播放
		this.play= function(){
			if(this.status=="stop"){

				hasPlay = true;
				//ctx.currentTime记录的为硬件的时间戳,非从播放时间开始计时,故播放前需要记录一个开始时间.
				this.startTime = this.audioContext.currentTime;
				console.log("startTime:"+this.startTime);

				this.status = "play";

				if(this.playStart >= this.buffer.duration){
					this.playStart = 0;
				}
				this.source.start(0,this.playStart);
				this.startPro();

			}else if(this.status="play"){
				this.stop();
			}else{
				console.log("状态异常");
			}
		};

		//开始进度条
		this.startPro= function(data,offset,samples){
			var self = this;
			this.interId = setInterval(function(){
				self.pro();
				//设置播放条的颜色
				$(".irs-bar,.irs-bar-edge").css({
					borderColor: self.progressColor,
					background: self.progressColor
				});
			},1000);
		};

		//停止进度条
		this.stopPro= function(){
			console.log("clear Inter");
			clearInterval(this.interId);
		};

		//进度条
		this.pro= function(){
			var self = this;
			var src = this.source;
			var buffer = this.buffer;
			var ctx = this.audioContext;
			console.log("hardware currentTimeStamp:"+ctx.currentTime);
			console.log("audio offset:" + this.playStart);
			console.log("audio startTime:"+this.startTime);
			console.log("played time:" + (parseFloat(ctx.currentTime) - parseFloat(this.startTime) + parseFloat(this.playStart)));
			var slider = $("#"+this.baseId+" .playSlider").data("ionRangeSlider");
			var from = parseFloat(ctx.currentTime) - parseFloat(this.startTime) + parseFloat(this.playStart);
			if(from >= buffer.duration){
				from = buffer.duration;
				this.playStart = buffer.duration;
				console.log("play end");
				this.stopPro();
				this.playBtn.src = this.playUrl;
				this.status = "stop";
                //重新初始化一个source
                this.source = this.audioContext.createBufferSource();
                this.gainNode = this.audioContext.createGain();
                this.source.buffer = this.buffer;
                this.source.connect(this.gainNode);
                this.gainNode.connect(this.audioContext.destination);
                this.gainNode.gain.value = this.volumeValue;
			}
			// Call sliders update method with any params
			slider.update({
				min: 0,
				max: buffer.duration,
				from:from,
				type: 'single',
				step: 0.1,
				postfix: " seconds",
				prettify: false,
				hasGrid: true,
				onFinish: function(obj){
					self.stop();
					self.playStart = parseFloat($("#"+self.baseId+" .playSlider").val());
					self.playBtn.src = self.pauseUrl;
					self.play();
				}
			});
			this.$curTime.text(this.formatTime(from));
		};

		this.formatTime=function(value) {
			var theTime = parseInt(value);// 秒
			var theTime1 = 0;// 分
			var theTime2 = 0;// 小时
			// alert(theTime);
			if(theTime > 60) {
				theTime1 = parseInt(theTime/60);
				theTime = parseInt(theTime%60);
				// alert(theTime1+"-"+theTime);
				if(theTime1 > 60) {
					theTime2 = parseInt(theTime1/60);
					theTime1 = parseInt(theTime1%60);
				}
			}
			var result;
			if(theTime < 10){
				result = "0"+ parseInt(theTime);
			}else{
				result = ""+ parseInt(theTime);
			}
			if(theTime1 >= 10) {
				result = ""+ parseInt(theTime1)+":"+result;
			}else if(theTime1 > 0 && theTime1 < 10){
				result = "0" + parseInt(theTime1) + ":" + result;
			}else{
				result = "00" + ":" + result;
			}
			if(theTime2 >= 10) {
				result = "" + parseInt(theTime2) + ":" + result;
			}else if(theTime2 > 0 && theTime2 < 10){
				result = "0" + parseInt(theTime2) + ":" + result;
			} else{
				result = "00" + ":" + result;
			}
			return result;
		};
	};

	return new Player();
});
