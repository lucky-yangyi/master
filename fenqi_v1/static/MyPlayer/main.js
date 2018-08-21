/**
 * Created by panfei on 2016/12/1.
 */
(function(baseUrl, initParams, urlArgs) {
    "use strict";

    requirejs.config({
        urlArgs: urlArgs,
        baseUrl: baseUrl,
        paths: {
            "jquery"    : "lib/jquery-2.1.4.min",
            "amrnb"  : "lib/amrnb",
            "ionRangeSlider"        : "js/ion.rangeSlider",
            "mplayer": "js/mplayer",
            "preloadimages": "js/preloadimages",
            "css": "lib/css.min",
            "ionrangeSlider" : "css/ion.rangeSlider",
            "ionrangeSlinderskinHTML5" : "css/ion.rangeSlider.skinHTML5"
        },
        config: {
            text: {
                useXhr: function (url, protocol, hostname, port) {
                    // 此处返回true，让text插件强制使用xhr加载，而不是通过jsonp加载
                    // （text资源位于不同域时，会默认使用jsonp，并加上.js获取）
                    // 参考：https://github.com/requirejs/text#xhr-restrictions
                    return true;
                }
            }
        }
    });

    require(["jquery", "mplayer", "preloadimages", "css!ionrangeSlider", "css!ionrangeSlinderskinHTML5"], function($, mplayer, preloadimages) {
        var $ = $.noConflict(true);

        var baseId = initParams.baseId;
        var playUrl = initParams.playUrl || (MyPlayer.getBaseUrl() + '/img/play.png');
        var pauseUrl = initParams.pauseUrl || (MyPlayer.getBaseUrl() + '/img/pause.png');
        var stopUrl = initParams.stopUrl || (MyPlayer.getBaseUrl() + '/img/stop.png');
        var playBarColor = initParams.playBarColor || '#F6F6F8';
        var progressColor = initParams.progressColor || '#808080';
        var audioUrl = initParams.audioUrl;

        //图片预加载，先加载图片再初始化录音组件
        preloadimages([playUrl, pauseUrl, stopUrl]).done(function() {
            mplayer.config(playUrl, pauseUrl, stopUrl, baseId, audioUrl, progressColor, playBarColor);
        });

    });

})(MyPlayer.getBaseUrl(), MyPlayer.getInitParams(), MyPlayer.urlArgs);
