/**
 * Created by panfei on 2016/12/1.
 */
(function(window) {
    "use strict";

    /**
     * 版本信息
     * @type {{name: string, date: string, svn: string}}
     */
    var version = {
        name: "MyPlayer",
        date: "2016-12-23",
        majorVersion: "4.5.4.3",
        git: "7d5a47a"
    };

    /**
     * 是否是开发模式（上线时须修改为false）
     */
    var devMode = true;

    /**
     * 服务器地址（支持http/https）
     * http默认端口80，https默认端口443
     */
    var config = {
        host: {
            "http:" : "http://10.130.41.151:9090"
        }
    };

    /**
     * 应用名称（默认为WA）
     */
    var appName = "MyPlayer";

    var MyPlayer = {};

    MyPlayer.devMode = devMode;

    MyPlayer.host = config.host[location.protocol];

    var baseUrl = MyPlayer.host + "/" + appName;

    MyPlayer.getBaseUrl = function() {
        return baseUrl;
    };

    /**
     * 加载模块附加的参数，开发时增加时间戳用于清除缓存，上线后去掉或设置固定值
     */
    MyPlayer.urlArgs = "v=" + (devMode ? new Date().getTime() : version.date + "_" + version.git);

    var initParams;

    MyPlayer.getInitParams = function() {
        return initParams;
    };

    /**
     * 判断当前浏览器是否支持
     * @returns {boolean}
     */
    function browserIsNotSupported() {
        return !(window.JSON && window.localStorage);
    }

    /**
     * 初始化入口函数
     * @param params
     */
    MyPlayer.init = function(params) {
        if (browserIsNotSupported()) {
            log("您的浏览器不支持录音组件，建议使用(谷歌/火狐)浏览器访问");
            return { code: -1, msg: "browser not supported."};
        }

        initParams = params || {};

        if (initParams.useLocal) {
            baseUrl = initParams.localUrl;
        }

        function loadMain() {
            loadScript((devMode ? "/main.js" : "/main.min.js") + "?" + MyPlayer.urlArgs);
        }

        if (isRequireJsExist()) {
            log("find requireJs exist. version: " + requirejs.version);
            loadMain();
        } else {
            loadScript("/lib/require-2.1.22.min.js", loadMain);
        }
    };

    /**
     * 判断页面是否已加载了requireJs库
     * @returns {requirejs|require|*|boolean}
     */
    function isRequireJsExist() {
        return window.requirejs && window.require && typeof window.requirejs.version === "string";
    }

    function log(msg) {
        if (window.console && window.console.log) {
            window.console.log(msg);
        }
    }

    function loadScript(url, callback) {
        var script = document.createElement("script");
        script.type = "text/javascript";
        script.async = true;
        script.src = baseUrl + url;
        script.onload = script.onreadystatechange = function() {
            if (!this.readyState || /complete|loaded/.test(this.readyState)) {
                log("script [" + url + "] load success.");
                if (typeof callback === "function") {
                    callback();
                }
                script.onload = null;
                script.onreadystatechange = null;
            }
        };
        (document.head || document.getElementsByTagName("head")[0]).appendChild(script);
    }

    if (!window.MyPlayer) {
        window.MyPlayer = MyPlayer;
    }

})(window);

