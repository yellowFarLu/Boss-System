+ function ($, Win) {
    "use strict";

    var GlobalFunc = new Global(),
        Id = location.href.split('?id=')[1],
        Tpl_tag_list = _.template($('#TplTagList').html().replace(/\&lt\;/ig, '\<')),
        Tpl_Comment_list = _.template($('#TplCommentlist').html().replace(/\&lt\;/ig, '\<')),
        $updateOpportForm = $('#updateOpportForm');

    /*-------------------------------------------------------------------------------------------
    ------------------------------------------------*/

    var Opport = function () {
        this.opportInfo = null;
        this.getOpportInfo();
        this.getAllComments();
        this.lookInfo = {
            page: 1,
            type: 0
        }
        this.EditOnClick = function (t) {
            this.saveAction(t);
            !t && this.changeOpportInfo();
        };
        this.saveAction = function (t) {
            $('#editBtn').attr("disabled", t);
            $('#saveBtn').attr("disabled", !t);

            $('#o_title').attr("disabled", !t);
            $('#sale').attr("disabled", !t);
            $('#Stage').attr("disabled", !t);
            $('#time').attr("disabled", !t);
            $('#Content').attr("disabled", !t);

            t ? $('#_tagAll').removeClass('hide') : $('#_tagAll').addClass('hide');
        }

        this.checkNewInfo = function (o) {
            o.title == this.opportInfo.title && delete o.title;
            o.quota == this.opportInfo.sales && delete o.quota;
            o.status == this.opportInfo.status && delete o.status;
            o.real_time == this.opportInfo.real_time && delete o.real_time;
            o.content == this.opportInfo.content && delete o.content;
            if (o.tags.toString() == this.opportInfo.tags.toString()) {
                o.tagsAlert = 0;
                delete o.tags;
            } else {
                o.tagsAlert = 1;
            }
            return o
        }
    }

    Opport.prototype = {
        getOpportInfo: function () {
            var _this = this;
            $.get('/getEmpBizInfo', {
                id: Id
            }, function (data) {
                var opport = data.biz;
                opport.tags = [];
                _.map(data.biz.Tags, function (m) {
                    opport.tags.push(m.tag_id);
                    var c = '.tag_' + m.tag_id;
                    $('#_tagAll').find(c).parent().addClass('on');
                });
                _this.opportInfo = opport;
                $('#o_title').val(data.biz.title);
                $('#o_customer').val(data.biz.customer.entName);
                $('#sale').val(data.biz.sales);
                $('#Stage option').each(function () {
                    if ($(this).val() == data.biz.status) {
                        $(this)[0].selected = true;
                    }
                })
                $('#estimate_time').val(data.biz.estimate_time);
                $('#time').val(data.biz.real_time);
                $('#Content').val(data.biz.content);

                var h = Tpl_tag_list(data.biz);
                $('#_updateTags').html(h);
            })
        },

        changeOpportInfo: function () {
            if (GlobalFunc.CheckForm($updateOpportForm)) return;
            var Tags = [];
            $('#_updateTags span').each(function () {
                Tags.push($(this).data('id'))
            });

            var opport = {
                opportunities_id: Id,
                title: $("#o_title").val(),
                quota: $('#sale').val(),
                status: $('#Stage option:checked').val(),
                real_time: $('#time').val(),
                content: $('#Content').val(),
                tags: Tags
            }

            opport = this.checkNewInfo(opport);
            if (Object.getOwnPropertyNames(opport).length == 2) {
                this.saveAction(false);
                return
            }
            $.post('/alert_opport', opport, function (data) {
                if (!!data && data.code == 0) {
                    GlobalFunc.Alert.success('更新成功 \^o^/');
                    opportCtl.getOpportInfo();
                    opportCtl.getAllComments();
                } else {
                    GlobalFunc.Alert.fail('更新失败');
                }
            })
        },

        addCommentBtnHtml: function () {
            var p = {
                id: Id,
                comment_content: $('#comment_area').val()
            }
            if (!$.trim(p.comment_content)) {
                GlobalFunc.Alert.info("备注不能为空!");
                return
            }
            $.post('/add_opport_comment', p, function (data) {
                if (!!data && data.code == 0) {
                    GlobalFunc.Alert.success('增加备注成功 \^o^/');
                    opportCtl.getOpportInfo();
                    opportCtl.getAllComments();
                } else {
                    GlobalFunc.Alert.fail('系统出错 X﹏X');
                }
            })
        },

        getAllComments: function (type, page, attr) {
            var p = {
                id: Id,
                type: type || 0,
                page: page || 1
            }
            attr = attr || 0;
            var _this = this;
            $.post("/getBizComments", p, function (data) {
                if (!!data && data.code == 0) {
                _.map(data.Comments, function (v) {
                    var _json = JSON.parse(v.comments)
                    v.comments = _json
                    v.text = ""
                    if (!v.type) {
                        v.text = _json["备注内容"]
                        return;
                    }
                    for (var i in _json) {
                        if (i == "action")
                            continue;
                        else
                            v.text += "[" + i + "] = " + _json[i] + "， "
                    }
                    v.text = v.text.substr(0, v.text.length - 2)
                });
                    var h = Tpl_Comment_list(data);
                    if (p.type == _this.lookInfo.type && attr != 0) {
                        $('#CommentListBox').append(h);
                    } else {
                        $('#CommentListBox').html(h);
                    }
                }
            });
            this.lookInfo = p
        },
        
        updateComment: function (id, v) {
            var der = $.Deferred()
            $.post('/alert_biz_comment', {
                biz_id: id,
                comments: v
            }, function (data) {
                if (!!data && data.code == 0) {
                    der.resolve(data)
                }
            })
            return der.promise()
        }
    }

    Win.opportCtl = new Opport();

    /*----------------------------------------------------------
    ------------------------------------------------*/


    $('body').on('click', '.recordBtn', function () {
        opportCtl.getAllComments(1);
    });

    $('body').on('click', '.lookMore', function () {
        opportCtl.getAllComments(opportCtl.lookInfo.type, opportCtl.lookInfo.page + 1, 1);
    });

    $('body').on('click', '.updateCommentBtn', function () {
        var textarea = $(this).parents('.panel').find('textarea');
        textarea.attr("readonly") == "readonly" ? textarea.removeAttr("readonly") : textarea.attr("readonly", true);
    });

    $('body').on('click', '.saveCommentBtn', function () {
        var dom = $(this).parents('.panel');
        var id = +dom.find('.comment_id').text();
        var value = dom.find('textarea').val()
        opportCtl.updateComment(id, value).then(function (data) {
            dom.find('textarea').attr("readonly", true);
        });
    });




}(jQuery, window)