"use strict";
(function ($) {

    $.fn.Sort = function (options) {

        //*-- 初始化参数
        var This = this;
        this.defaults = {
            btnId: 'sureBtn',
            sortType: 'default',
            sortKey: 'timex',
            valueType: 'time',
            
            tag_html: "<%for(var i=0;i<list.length;i++){%>" +
                "<span data-id='<%=list[i].tag_id%>' class='tag'><%=list[i].Name%></span>" +
                "<%}%>",
            
            area_html: "<%for(var i=0;i<list.length;i++){%>" +
                "<span data-id='<%=list[i].province_id%>' class='tag'><%=list[i].province_name%></span>" +
                "<%}%>",

            html: '<section class="panel sortTip"> ' +
                ' <header class="panel-heading">' +
                '条件筛选<a class="fa fa-times pull-right" href="javascript:;"></a>' +
                '</header>' +
                '<div class="panel-body">' +
                '<form role="form" data-key="<%sortKey%>">' +
                '<div class="row">' +
                '<div class="col-md-12">' +
                '<%sort%>' +
                '<%type%>' +
                '</div>' +
                '</div>' +
                '<button type="button" class="btn btn-success pull-right sortGo" id="<%sureBtn%>" data-key="<%sortKey%>">确定</button>' +
                '</form>' +
                '</div>' +
                '</section>' +
                '<div class="mask"></div>',

            type: {
                time: '<div class="form-group valueLine">' +
                    '<label for="">筛选值：</label>' +
                    '<div class="input-group input-large custom-date-range">' +
                    '<input type="text" class="form-control startValue default-date-picker" name="from">' +
                    '<span class="input-group-addon">To</span>' +
                    '<input type="text" class="form-control endValue default-date-picker" name="to">' +
                    '</div>' +
                    '</div>',

                tag: '<div class="form-group valueLine">' +
                    '<label for="">筛选值：</label>' +
                    '<div id="allTagsBox" class="tagsinput">' +
                    '</div>' +
                    '</div>',
                
                area: '<div class="form-group valueLine">' +
                    '<label for="">筛选值：</label>' +
                    '<div id="allAreaBox" class="tagsinput">' +
                    '</div>' +
                    '</div>'
            },

            sort: {
                none: '',
                default: '<div class="form-group orderBy">' +
                    '<label for="">排序：</label>' +
                    '<div class="row icheck ">' +
                    '<div class="square-purple col-sm-6">' +
                    '<div class="radio ">' +
                    '<input data-index="1" type="radio"  name="demo-radio" >' +
                    '<label>正序</label>' +
                    '</div>' +
                    '</div>' +
                    '<div class="square-purple col-sm-6">' +
                    '<div class="radio ">' +
                    '<input data-index="0" type="radio"  name="demo-radio" checked>' +
                    '<label>倒序</label>' +
                    '</div>' +
                    '</div>' +
                    '</div>' +
                    '</div>'
            }
        };

        this.opts = $.extend({}, this.defaults, options);

        //*-- 事件绑定
        this.eventBind = function () {
            var h = this.opts.html.replace(/<%sureBtn%>/ig, this.opts.btnId)
                .replace(/<%sortKey%>/ig, this.opts.sortKey)
                .replace(/<%sort%>/ig, this.opts.sort[this.opts.sortType])
                .replace(/<%type%>/ig, this.opts.type[this.opts.valueType]);

            this.append(h);

            var btn = this.siblings('.sortBtn');

            btn.click(function () {
                $('.sortBox').fadeOut(300);
                $(this).siblings('.sortBox').fadeToggle(300);
            });

            $('body').on('click', '.sortBox .panel-heading a,.sortGo', function () {
                $(this).parents('.sortBox').fadeOut(300);
            });
            
            $('.mask').click(function (e) {
                e.stopPropagation();
                $('.sortBox').fadeOut(300);
            })
        }

        this.FormInit = function (form) {
            var _form = form ? form : this.find('form');
            _form[0].reset();
        }

        this.GetAllTag = function () {
            $('body').on('click', '#allTagsBox .tag', function () {
                $(this).toggleClass('on');
            });
            $.post('/get_all_tag', function (data) {
                if (!!data && data.code == 0) {
                    var h = _.template(This.opts.tag_html)(data);
                    $('#allTagsBox').html(h);
                }
            })
        }

        this.GetAllArea = function () {
            $('body').on('click', '#allAreaBox .tag', function () {
                $(this).toggleClass('on');
            });
            $.post('/get_all_province', function (data) {
                if (!!data && data.code == 0) {
                    var h = _.template(This.opts.area_html)(data);
                    $('#allAreaBox').html(h);
                }
            })
        }


        //*-- 初始化函数
        this.init = function () {
            This.eventBind();
            This.opts.valueType == "tag" && This.GetAllTag();
            This.opts.valueType == "area" && This.GetAllArea();
        }();
        return this
    }
})(jQuery)