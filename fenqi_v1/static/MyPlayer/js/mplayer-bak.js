var MyPlayer = function(){
		this.baseUrl = null;
		this.baseId = null;
		this.audioUrl = null;
		this.config = function(url,audioUrl,id){
			this.baseUrl = url;
			this.baseId = id;
			this.audioUrl = audioUrl;
			this.genPlayer();
		},
			this.genPlayer = function(){
				$("#"+this.baseId).empty();
				$("#"+this.baseId).append(
					'<div class="sample-amr box-shadow" style="background:#F6F6F8;width:100%;height: 47px;position: relative;display: -webkit-flex;display: flex;line-height: 47px;min-width: 680px;">'+
						'<div style="padding:4px 32px;display: inline-block;vertical-align: middle;flex-shrink: 0;">'+
							//'<span class="fa fa-play-circle play"></span>'+
							//'<span class="fa fa-pause-circle stop"></span>'+
							'<img class="play" src="img/play1.png" width="37px" height="37px" style="margin-right: 15px;">'+
							'<img class="stop" src="img/stop.png" width="37px" height="37px">'+
						'</div>'+
						'<span class="curTime">00:00:00</span>'+
						'<div style="flex-shrink: 1;flex-grow: 2;">'+
							'<input  class="playSlider" type="hidden" name="playSlider" value="" />'+
						'</div>'+
						'<span class="totalTime">00:00:00</span>'+
						'<img class="sound" src="img/sound.png">'+
						'<div class="v" style="width:80px;display: inline-block;flex-shrink: 0;padding: 17px 20px 17px 0px;">'+
							'<input  class="volumeSlider" type="hidden" name="volumeSlider" value="" />'+
						'</div>'+
					'</div>'
				);
				this.initEventHandler();
			},
			this.initEventHandler = function(){
				var mp = this;
				//播放条
				console.log($("#"+mp.baseId+" .playSlider"));
				$("#"+mp.baseId+" .playSlider").ionRangeSlider({
					type: 'single',
					min:0,
					hide_min_max: true,
					hide_from_to: true,
					onChange:function(obj){
						//播放时拖动进度条
						if(boxes.status=="play"){
							boxes.stop();
						}
					}
				});
				var boxes = new Boxes();
				boxes.baseUrl = mp.baseUrl;
				boxes.baseId = mp.baseId;
				//音量条
				$("#"+this.baseId+" .volumeSlider").ionRangeSlider({
					type: 'single',
					min:0,
					max:1,
					step: 0.1,
					from:boxes.volumeValue,
					hide_min_max: true,
					hide_from_to: true,
					force_edges:true,
					onFinish: function(obj){
						boxes.volumeValue = parseFloat($("#"+boxes.baseId+" .volumeSlider").val());
						if(boxes.gainNode){
							boxes.gainNode.gain.value = boxes.volumeValue;
						}
					}
				});
				boxes.initAPI();
				boxes.$curTime = $("#"+boxes.baseId+" .sample-amr .curTime");
				boxes.$totalTime = $("#"+boxes.baseId+" .sample-amr .totalTime");
				//播放暂停
				E("#"+mp.baseId+" .sample-amr .play").onclick = function() {
					//保存按钮引用
					boxes.playBtn = this;
					if(boxes.status=="stop"){
						this.src=boxes.baseUrl+"img/pause1.png";
						boxes.fetchBlob(mp.audioUrl, function(blob) {
							//这之内的数据解码和加载都是需要时间的,故播放操作要写在其内的回调中
							boxes.playAmrBlob(blob);
						});
					}else if(boxes.status=="play"){
						this.src= boxes.baseUrl+"img/play1.png";
						boxes.playStart = parseFloat($("#"+boxes.baseId+" .playSlider").val());
						boxes.stop();
					}else{
						alert("状态异常");
					}

				};
				//终止
				E("#"+mp.baseId+" .sample-amr .stop").onclick = function() {
					boxes.playBtn.src=boxes.baseUrl+"/img/play.png";
					boxes.terminal();
				}
			}

	}




	function E(selector) {
		return document.querySelector(selector);
	}

	var Boxes = function() {
		this.baseId = null;
		this.baseUrl = null;
		this.playBtn = null;     //播放按钮引用
		this.interId = null;     //定时调度id
		this.status = "stop"; //当前播放状态
		this.startTime = 0;  //开始播放时间
		this.totalTime = 0;  //总时间
		this.$curTime = null;  //显示当前时间jquery对象
		this.$totalTime = null;  //显示总时间的jquery对象
		this.playStart = 0;  //音频文件起始播放时间点
		this.buffer = null; //要处理的文件
		this.audioContext = null; //进行音频处理的上下文，稍后会进行初始化
		this.source = null; //保存音频
		this.gainNode = null; //音频处理节点
		this.volumeValue = 0.5; //音量值

		this.initAPI= function() {
			//统一前缀，方便调用
			window.AudioContext = window.AudioContext || window.webkitAudioContext || window.mozAudioContext || window.msAudioContext;
			//安全地实例化一个AudioContext并赋值到audioContext属性上，方便后面处理音频使用
			try {
				this.audioContext = new AudioContext();
			} catch (e) {
				console.log('!您的浏览器不支持AudioContext:(');
				console.log(e);
			}
		},
			this.fetchBlob= function(url, callback) {
				var xhr = new XMLHttpRequest();
				xhr.open('GET', url);
				xhr.responseType = 'blob';
				xhr.onload = function() {
					callback(this.response);
				};
				xhr.onerror = function() {
					alert('Failed to fetch ' + url);
				};
				xhr.send();
			},
			this.readBlob= function(blob, callback) {
				var reader = new FileReader();
				reader.onload = function(e) {
					var data = new Uint8Array(e.target.result);
					callback(data);
				};
				reader.readAsArrayBuffer(blob);
			},
			this.playAmrBlob= function(blob, callback) {
				var self = this;
				this.readBlob(blob, function(data) {
					self.playAmrArray(data);
				});
			},
			this.playAmrArray= function(array) {
				var samples = AMR.decode(array);
				if (!samples) {
					alert('Failed to decode!');
					return;
				}
				this.playPcm(samples);
			},
			this.playPcm= function(samples) {
				this.source = this.audioContext.createBufferSource();
				this.gainNode = this.audioContext.createGain();
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
				this.source.buffer = this.buffer;
				this.source.connect(this.gainNode);
				this.gainNode.connect(this.audioContext.destination);
				this.gainNode.gain.value = this.volumeValue;
				this.totalTime = this.formatTime(this.buffer.duration);
				console.log("totalTime:"+this.totalTime)
				this.$totalTime.text(this.totalTime);
				this.play();
			},
			//停止(暂停)
			this.stop= function(){
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
			},
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
			},
			//播放
			this.play= function(){
				if(this.status=="stop"){
					//ctx.currentTime记录的为硬件的时间戳,非从播放时间开始计时,故播放前需要记录一个开始时间.
					this.startTime=this.audioContext.currentTime;
					console.log("startTime:"+this.startTime);

					this.status="play";

					if(this.playStart>=this.buffer.duration){
						this.playStart = 0;
					}
					this.source.start(0,this.playStart);
					this.startPro();

				}else if(this.status="play"){
					this.stop();
				}else{
					console.log("状态异常");
				}
			},
			//开始进度条
			this.startPro= function(data,offset,samples){
				var self = this;
				this.interId = setInterval(function(){self.pro();},1000);
			},
			//停止进度条
			this.stopPro= function(){
				console.log("clear Inter");
				clearInterval(this.interId);
			},
			//进度条
			this.pro= function(){
				var self = this;
				var src = this.source;
				var buffer = this.buffer;
				var ctx = this.audioContext;
				console.log("hardware currentTimeStamp:"+ctx.currentTime);
				console.log("audio offset:"+this.playStart);
				console.log("audio startTime:"+this.startTime);
				console.log("played time:"+(parseFloat(ctx.currentTime)-parseFloat(this.startTime)+parseFloat(this.playStart)));
				var slider = $("#"+this.baseId+" .playSlider").data("ionRangeSlider");
				var from = parseFloat(ctx.currentTime)-parseFloat(this.startTime)+parseFloat(this.playStart);
				if(from>=buffer.duration){
					from = buffer.duration;
					this.playStart = buffer.duration;
					console.log("play end");
					this.stopPro();
					this.playBtn.src=this.baseUrl+"/img/play.png";
					this.status = "stop";
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
						self.playBtn.src=self.baseUrl+"img/pause1.png";
						self.play();
					}
				});
				this.$curTime.text(this.formatTime(from));
			},
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
				if(theTime<10){
					result = "0"+parseInt(theTime);
				}else{
					result = ""+parseInt(theTime);
				}
				if(theTime1 >= 10) {
					result = ""+parseInt(theTime1)+":"+result;
				}else if(theTime1 >0&&theTime1<10){
					result = "0"+parseInt(theTime1)+":"+result;
				}else{
					result = "00"+":"+result;
				}
				if(theTime2 >= 10) {
					result = ""+parseInt(theTime2)+":"+result;
				}else if(theTime2 >0&&theTime2<10){
					result = "0"+parseInt(theTime2)+":"+result;
				} else{
					result = "00"+":"+result;
				}
				return result;
			}
	};
