/**
 * Created by panfei on 2016/12/8.
 */
define([], function() {
    //图片预加载
    function preloadimages(arr) {
        var images = [];
        var loadedimages = 0;
        var postaction = function() {}; //此处增加一个postaction函数

        function imageloadpost() {
            loadedimages++;
            if(loadedimages === arr.length) {
                console.log("图片加载成功");
                postaction(images); // 加载完成后调用postaction函数并将image数组作为参数传递进去
            }
        }
        for(var i = 0; i < arr.length; i++) {
            images[i] = new Image();
            images[i].onload = function() {
                imageloadpost();
            };
            images[i].onerror = function() {
                imageloadpost();
            }
            images[i].src = arr[i];
        }

        return{  //此处返回一个空白对象的done方法
            done: function(func) {
                postaction = func || postaction
            }
        }
    }

    return preloadimages;
});