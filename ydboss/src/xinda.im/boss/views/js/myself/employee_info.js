+ function ($, Win) {
    "use strict";

    var GlobalFunc = new Global(),
        Id = location.href.split('?id=')[1],
        Tpl_comment_list = _.template($('#TplComment').html().replace(/\&lt\;/ig, '\<')),
        $customCommentBox = $('#customCommentBox'),
        $opportCommentBox = $('#opportCommentBox'),
        $cLookMore = $('#cLookMore'),
        $oLookMore = $('#oLookMore')

    /* ----------------------------------------------------------------------------------------------- */
    var AccountInfo = function () {
        this.page = {
            custom: 0,
            opport: 0
        };
        this.getComments("custom");
        this.getComments("opport");
    }
    AccountInfo.prototype = {
        getComments: function (type, page) {
            var This = this;
            var page = page || 1;
            var query = type == "custom" ? "/get_custom_comments" : "/get_opport_comments";
            var dom = type == "custom" ? $customCommentBox : $opportCommentBox;
            $.get(query, {
                id: Id,
                page: page
            }, function (data) {
                if (!!data && data.code == 0) {
                    var h = Tpl_comment_list(data).replace(/<\/br>/ig, "").replace(/:/ig, ": ");
                    page == 1 ? dom.html(h) : dom.append(h);
                    if (type == "custom") {
                        This.page.custom++;
                        data.list.length == 0 && $cLookMore.parent().hide();
                    } else {
                        This.page.opport++;
                        data.list.length == 0 && $oLookMore.parent().hide();
                    }
                }
            })
        }
    }
    var AccInfoCtl = new AccountInfo();
    /* ----------------------------------------------------------------------------------------------- */


    $cLookMore.click(function () {
        var page = AccInfoCtl.page.custom + 1;
        AccInfoCtl.getComments("custom", page);
    });
    $oLookMore.click(function () {
        var page = AccInfoCtl.page.opport + 1;
        AccInfoCtl.getComments("opport", page);
    });

}(jQuery, window)