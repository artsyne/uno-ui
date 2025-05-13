(function($) {
    // 定义插件
    $.fn.ajaxLoader = function(options) {
        // 默认设置
        var settings = $.extend({
            contentSelector: '#content', // 内容插入的选择器
            onError: function(error) { // 默认错误处理
                console.error('Error:', error);
            }
        }, options);

        // 监听 <a role="ajax"> 链接的点击事件
        this.on('click', 'a[role="ajax"]', function(event) {
            event.preventDefault(); // 阻止默认的链接跳转行为

            var url = $(this).attr('href');

            // 使用 jQuery 的 $.ajax 方法发送 AJAX 请求
            $.ajax({
                url: url,
                method: 'GET',
                success: function(data) {
                    $(settings.contentSelector).html(data); // 将返回的内容插入到指定的部分页面中
                    history.pushState(null, '', url); // 更新浏览器地址
                },
                error: function(jqXHR, textStatus, errorThrown) {
                    settings.onError(textStatus, errorThrown);
                }
            });
        });

        // 处理浏览器的前进和后退按钮
        window.onpopstate = function() {
            var url = location.pathname;
            $.ajax({
                url: url,
                method: 'GET',
                success: function(data) {
                    $(settings.contentSelector).html(data); // 将返回的内容插入到指定的部分页面中
                },
                error: function(jqXHR, textStatus, errorThrown) {
                    settings.onError(textStatus, errorThrown);
                }
            });
        };

        return this;
    };
}(jQuery));