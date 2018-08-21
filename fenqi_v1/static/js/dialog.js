/* 
    author : liaojia
    date : 2017-12-13
    version : 2017-12-13
*/
var dialog = {
    default: {
        header :{
            style : "green",
        },
        close : {
            style : "green"
        },
        title: "提示",//弹窗标题 选填 默认为提示
        select: {
            html: "",
            index: 1
        },//如需select  传入select代码片段 选填
        textarea: {
            html: "",
            index: 1
        },//如需textarea  传入textarea代码片段 选填
        input: {
            html: "",//由于input种类过多， 直接传入代码片段
            index: 1
        },//如需input  传入input代码片段 选填
        others: {
            html: ""
        },
        tips: "",//提示文案
        buttons: [{
            text: "确定",
            style: "green",
            func: function () {

            }
        }]
    },
    options: {},
    init: function (options) {
        this.options = options;
        this.render();
        this.bindEvent();
    },
    render: function () {
        var dialogHtml = [];
        dialogHtml.push('<div class="g-modal J_modal">');
        dialogHtml.push('<div class="g-content">');
        dialogHtml.push('<div class="modal-title clearfix ' + (this.options.header&&this.options.header.style || this.default.header.style) + '">');
        dialogHtml.push('<div class="title-text" id="payplan-title">' + (this.options.title || this.default.title) + '</div>');
        dialogHtml.push('<div class="modal-close ' + (this.options.close&&this.options.close.style || this.default.close.style) + '">&times;</div>');
        dialogHtml.push('</div>');
        dialogHtml.push('<div class="modal-body">');
        if (this.options.tips) {
            dialogHtml.push('<div class="modal-tips" id="callingPhone">' + this.options.tips + '</div>');
        }
        if (this.options.input) {
            dialogHtml.push(this.options.input.html);
        }
        if (this.options.select) {
            dialogHtml.push(this.options.select.html);
        }
        if (this.options.textarea) {
            dialogHtml.push(this.options.textarea.html);
        }
        if (this.options.others) {
            dialogHtml.push(this.options.others.html);
        } else {
            dialogHtml.push('<div class="red-text"></div>');
        }
        dialogHtml.push('</div>');
        dialogHtml.push('<div class="modal-foot">');
        var len = this.options.buttons && this.options.buttons.length || this.default.buttons.length;
        for (var i = 0; i < len; i++) {
            var item = this.options.buttons && this.options.buttons[i] || this.default.buttons[i]
            dialogHtml.push('<div data-index='+i+' class="btn btn-' + item.style + '">' + item.text + '</div>');
        }
        dialogHtml.push('</div>');
        dialogHtml.push('</div>');
        dialogHtml.push('</div>');
        $("body").append(dialogHtml.join(""));
    },
    bindEvent: function () {
        var me = this;
        $(".modal-close").off().click(function () {//关闭按钮
            $(".J_modal").remove();
        });
        var len = $(".J_modal .modal-foot .btn").length;
        for (var i = 0; i < len; i++) {
            $(".J_modal .modal-foot .btn").eq(i).off().click(function (e) {
                var index = parseInt($(this).data("index"));
                if(me.options.buttons&&me.options.buttons[index].func) {
                    var result = me.options.buttons[index].func(e);
                    if(!result){
                        return;
                    }
                }else{
                    $(".J_modal").remove();
                }
            })
        }
    }
}
// demo
// dialog.init({
//     title: "测试",
//     textarea: {
//         placeholder: "测试"
//     },
//     input: {
//         html: '<input id="phoneTime"  class="Wdate form-control" value="" placeholder="入列时间" style="margin-bottom:10px;" type="text" onFocus="WdatePicker({dateFmt:&quot;yyyy-MM-dd HH:mm:ss&quot;,minDate:&quot;%y-%M-%d&quot;})"/>'
//     },
//     buttons: [{
//         text: "确定",
//         style: "green"
//     }, {
//         text: "取消",
//         style: "green"
//     }]
// })