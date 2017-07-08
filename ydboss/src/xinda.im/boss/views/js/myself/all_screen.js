/**
 * @params {$: jQuery}
 * @params {WIN: window}
 * @params {DOM: $('body')}
 */
+ function ($, WIN, DOM) {
    "use strict"
    
    // Statements
    var GlobalFunc = new Global(),
        $provinceList = $('#provinceList'),
        $tagList = $('#tagList'),
        $provinceInput = $('#provinceInput'),
        $tagInput = $('#tagInput'),
        $provincetagList = $('#provincetagList'),
        $formData = $('#formData'),
        $clearBtn = $('#clearBtn'),
        $provincetagList = $('#provincetagList'),
        $TagtagList = $('#TagtagList'),
        $allProvinceList = $('#allProvinceList'),
        $allTagsList = $('#allTagsList'),
        $provinceName = $('#provinceName'),
        $tagName = $('#tagName'),
        $provinceSearchAll = $('#provinceSearchAll'),
        $tagsSearchAll = $('#tagsSearchAll'),
        $ModalProvince = $('#ModalProvince'),
        $ModalTags = $('#ModalTags'),
        $chooseProvinceBtn = $('#chooseProvinceBtn'),
        $chooseTagBtn = $('#chooseTagBtn')

    /*
     # CustomScreen [CLASS]
     */
    function Screen() {
        this.provinces = [];
        this.tags = [];
        
        this.TplProvinceHtml = _.template('<%for(var i=0;i<list.length;i++){%>' +
            '<div data-id="<%=list[i].province_id%>" class="screen_li"><i class="fa fa-plus-square-o"></i><span class="name"><%=list[i].province_name%></span></div>' +
            '<%}%>');
        
        this.TplTagsHtml = _.template('<%for(var i=0;i<list.length;i++){%>' +
            '<div data-id="<%=list[i].tag_id%>" class="screen_li"><i class="fa fa-plus-square-o"></i><span class="name"><%=list[i].Name%></span></div>' +
            '<%}%>');
        
        this.TplOneTagHtml = _.template('<span class="tag_li" data-id="<%=id%>"><i class="fa fa-tag"></i> <%=name%><span class="close">x</span></span>');
        
        this.TplProvinceSearchHtml = _.template('<%for(var i=0;i<list.length;i++){%>' +
            '<div data-id="<%=list[i].province_id%>" class="all_screen_li province_<%=list[i].province_id%>"><i class="fa fa-plus-square"></i><span class="name"><%=list[i].province_name%></span><button class="btn btn-sm chooseBtn pull-right" type="button">选择</button></div>' +
            '<%}%>');
        
        this.TplTagsSearchHtml = _.template('<%for(var i=0;i<list.length;i++){%>' +
            '<div data-id="<%=list[i].tag_id%>" class="all_screen_li tag_<%=list[i].tag_id%>"><i class="fa fa-plus-square"></i><span class="name"><%=list[i].Name%></span><button class="btn btn-sm chooseBtn pull-right" type="button">选择</button></div>' +
            '<%}%>');
        
        // check province choose add class 'on'
        this._checkProvinceChoose = function () {
            _.map(this.provinces, function (v) {
                var _class = ".province_" + v;
                $allProvinceList.find(_class).addClass('on');
            })
        };
        
        // check tag choose add class 'on'
        this._checkTagChoose = function () {
            _.map(this.tags, function (v) {
                var _class = ".province_" + v;
                $allTagsList.find(_class).addClass('on');
            })
        }

    };

    /*
     # CustomScreen Prototype {Http Api}
     */
    Screen.prototype = {        
        
        // get all province
        getAllProvince: function () {
            $.post('/get_all_province', function (data) {
                if (!!data && data.code == 0) {
                    var h = ScreenCtl.TplProvinceSearchHtml(data);
                    $allProvinceList.html(h);
                    ScreenCtl._checkProvinceChoose();
                }
            })
        },
        
        // get province by key
        getProvinceByKey: function (key) {
            var der = $.Deferred();
            $.post('/get_province_by_key', {
                key: key
            }, function (data) {
                if (!!data && data.code == 0) {
                    der.resolve(data)
                }
            });
            return der.promise()
        },
        
        // get all tags
        getAllTags: function () {
            $.post("/get_all_tag", function (data) {
                if (!!data && data.code == 0) {
                    var h = ScreenCtl.TplTagsSearchHtml(data);
                    $allTagsList.html(h);
                    ScreenCtl._checkTagChoose();
                }
            })
        },

        // get tag by key
        getTagByKey: function (key) {
            var der = $.Deferred();
            $.post('/get_key_tag', {
                key: key
            }, function (data) {
                if (!!data && data.code == 0) {
                    der.resolve(data)
                }
            });
            return der.promise()
        }
    };

    var ScreenCtl = new Screen();

    /*
     # Event Bind
     */
    GlobalFunc.Focus('#provinceSearchAll', '#provinceForm');
    GlobalFunc.Focus('#tagsSearchAll', '#tagsForm');
    ScreenCtl.getAllProvince();
    ScreenCtl.getAllTags();

    DOM.on('keyup focus', '#provinceInput', function () {
        var $this = $(this);
        var h = $this.parents('.screenStage').height() + 4 + 'px';
        var w = $this.parents('.screenStage').width()
        $provinceList.html("").css("top", h).show();
        $this.parents('.form-group.screen').siblings('.screen').find('.focusShow').hide();
        ScreenCtl.getProvinceByKey($this.val()).then(function (data) {
            var h = ScreenCtl.TplProvinceHtml(data);
            $provinceList.html(h);
        });
    });

    DOM.on('click', '.focusShow .screen_li', function () {
        var $this = $(this);
        var dom = $this.parent();
        var o = {
            id: $this.data('id'),
            name: $this.children('.name').text()
        };
        dom.hide();
        var element = {
            "provinceList": "provinces"
        }[dom.attr("id")] || "tags";
        if (_.includes(ScreenCtl[element], o.id)) {
            return
        } else {
            ScreenCtl[element].push(o.id);
        }
        var h = ScreenCtl.TplOneTagHtml(o);
        dom.siblings('.screenStage').find('.tagStage').append(h).siblings('.keyInput').val("");
    });

    DOM.on('click', '.tag_li .close', function () {
        var $this = $(this);
        var DOM = $this.parents('.tagStage').first();
        var id = $this.parent('.tag_li').data('id');
        var element = {
            "provincetagList": "provinces"
        }[DOM.attr("id")] || "tags";
        _.remove(ScreenCtl[element], function (v) {
            return v == id
        });
        $this.parent('.tag_li').remove();
    });

    DOM.on('click', '.screenStage', function () {
        $(this).find('.keyInput').focus();
    });

    DOM.on('click', '*:not("td")', function (e) {
        e.stopPropagation();
        if ($(this).hasClass('screenStage') || $(this).parent().hasClass('screenStage') || $(this).hasClass('screen_li')) return;
        $('.focusShow ').hide();
    });

    DOM.on('keyup focus', '#tagInput', function () {
        var $this = $(this);
        var h = $this.parents('.screenStage').height() + 4 + 'px';
        $tagList.html("").css("top", h).show();
        ScreenCtl.getTagByKey($this.val()).then(function (data) {
            var h = ScreenCtl.TplTagsHtml(data);
            $tagList.html(h);
        });
    });

    /*
     # Clear All Query Options
     */

    $clearBtn.click(function () {
        $formData[0].reset();
        $provincetagList.empty();
        $TagtagList.empty();
        ScreenCtl.provinces = ScreenCtl.tags = [];
    });

    /*
     # Province Search
     */
    $provinceName.keyup(function () {
        var v = $.trim($(this).val());
        if (!v) {
            ScreenCtl.getAllProvince();
            return;
        }
        ScreenCtl.getProvinceByKey(v).then(function (data) {
            var h = ScreenCtl.TplProvinceSearchHtml(data);
            $allProvinceList.html(h);
            ScreenCtl._checkProvinceChoose();
        })
    });

    $provinceSearchAll.click(function () {
        $ModalProvince.modal("show");
        $allProvinceList.children('.all_screen_li').removeClass('on');
        ScreenCtl._checkProvinceChoose();
    });

    $chooseProvinceBtn.click(function () {
        var h = "",
            provinces = [];
        $allProvinceList.find('.all_screen_li.on').each(function () {
            var o = {
                id: $(this).data("id"),
                name: $(this).find(".name").text()
            };
            provinces.push(o.id);
            h += ScreenCtl.TplOneTagHtml(o);
        });
        $provincetagList.empty().html(h);
        ScreenCtl.provinces = provinces;
    });

    /*
     # Tag Search
     */
    $tagName.keyup(function () {
        var v = $.trim($(this).val());
        if (!v) {
            ScreenCtl.getAllTags();
            return;
        }
        ScreenCtl.getTagByKey(v).then(function (data) {
            var h = ScreenCtl.TplTagsSearchHtml(data);
            $allTagsList.html(h);
            ScreenCtl._checkTagChoose();
        })
    });

    $tagsSearchAll.click(function () {
        $ModalTags.modal("show");
        $allTagsList.children('.all_screen_li').removeClass('on');
        ScreenCtl._checkTagChoose();
    });

    $chooseTagBtn.click(function () {
        var h = "",
            tags = [];
        $allTagsList.find('.all_screen_li.on').each(function () {
            var o = {
                id: $(this).data("id"),
                name: $(this).find(".name").text()
            };
            tags.push(o.id);
            h += ScreenCtl.TplOneTagHtml(o);
        });
        $TagtagList.empty().html(h);
        ScreenCtl.tags = tags;
    });

    /*
     # Common
     */
    DOM.on('click', '.all_screen_li', function () {
        $(this).toggleClass('on');
    });

    DOM.on('click', '.chooseBtn', function () {
        $(this).parent().toggleClass('on');
    });

}(jQuery, window, $("body"))