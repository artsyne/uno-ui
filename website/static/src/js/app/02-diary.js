/**
 * 日记本管理插件
 * 
 * 此插件提供日记本管理相关功能，包括创建新日记本、
 * 表单交互、主题选择等功能。
 * 
 * @author nobody
 * @version 1.0.0
 */
(($) => {
    $.fn.diary = function (options) {
      // 默认设置与用户设置合并
      const settings = $.extend({
        ajaxUrl: '',         // API请求地址
        debug: true,        // 是否启用调试模式
        // 事件绑定配置，便于集中管理所有按钮事件
        eventBindings: [
          {
            selector: '[jq-bind-click="addDaybook"]',
            event: 'click',
            handler: 'showDiaryBox'
          },
          {
            selector: '#post_btn',
            event: 'click',
            handler: 'postInfo'
          }
        ]
      }, options);
  
      // 调试日志函数，仅在debug模式开启时输出
      const log = (message, ...data) => {
        if (settings.debug) {
          console.log(`[Diary] ${message}`, ...data);
        }
      };
  
      // 日记本相关方法
      const methods = {
        /**
         * 显示新增日记本对话框
         * @param {Event} e - 事件对象
         */
        showDiaryBox(e) {
          log('显示新增日记本对话框');
          if (e) e.preventDefault();
          
          $.box.dialog({
            title: '新增日记本',
            content: '<div class="modal-body">正在加载中......</div>',
            width: 700,
            height: 400
          });
        },
  
        /**
         * 模态框显示后初始化表单
         * 注：将来会替换为自定义实现
         */
        initDaybookForm() {
          log('初始化日记本表单');
          
          $('#daybook_add_modal').on('shown.bs.modal', () => {
            // 加载表单HTML
            $.get(`${settings.ajaxUrl}/form/daybook`, (data) => {
              $('#daybook_add_modal .form_html:first').html(data);
            }).done(() => {
              // 初始化表单验证
              if ($.fn.valid) {
                $('#daybook_editor_form').valid();
              }
              
              // 绑定主题选择事件
              $('.d-theme>div')
                .mouseover(function () {
                  $(this).addClass('dtfhover');
                })
                .mouseout(function () {
                  $(this).removeClass('dtfhover');
                })
                .click(function () {
                  $('.d-theme>div').removeClass('dtfselected');
                  $(this).addClass('dtfselected');
                  $('#theme_id').val($(this).attr('theme'));
                });
            });
          });
  
          // 模态框关闭时清理表单
          $('#daybook_add_modal').on('hidden.bs.modal', () => {
            $('#daybook_add_modal .form_html:first').html(
              '<div class="modal-body">正在加载中......</div>'
            );
          });
        },
  
        /**
         * 提交日记本表单数据
         * 创建新日记本
         */
        doDaybook() {
          log('提交日记本表单');
          
          // 表单验证
          const form = $('#daybook_editor_form');
          if ($.fn.valid && form.valid('checkError')) {
            return;
          }
          
          // 提交表单数据
          $.post(
            `${settings.ajaxUrl}/daybook`,
            form.serialize(),
            (data) => {
              if (data.Status === 0) {
                // 创建成功，添加到列表
                const html = `
                  <div class="col-sm-6 col-md-4 col-lg-3">
                    <div class="card card-diary-theme dt${data.Data.theme_id}">
                      <div class="card-body">
                        <h5>
                          ${data.Data.view_setting === 2 ? "<i class='mdi mdi-lock-outline'></i>" : ''}
                          <a href='/d/${data.Data.diary_id}'>${$.HtmlEncode(data.Data.title)}</a>
                        </h5>
                        <p>0篇</p>
                      </div>
                    </div>
                  </div>`;
                
                $('#daybook_list').prepend(html);
                $('#daybook_add_modal').modal('hide');
                
                log('日记本创建成功', data.Data);
              } else {
                log('日记本创建失败', data);
              }
            },
            'json'
          ).fail((xhr) => {
            log('日记本创建请求失败', xhr.statusText);
          });
        },
  
        /**
         * 提交信息表单
         * @param {Event} e - 事件对象
         */
        postInfo(e) {
          log('提交信息表单');
          if (e) e.preventDefault();
          
          const form = $(this).closest('form');
          
          // 表单验证
          if ($.fn.valid && form.valid('checkError')) {
            return;
          }
          
          // 禁用表单控件，防止重复提交
          if ($.utils && $.utils.disableFormInputs) {
            $.utils.disableFormInputs(form, true);
          }
          
          // 提交表单数据
          $.post(
            `${settings.ajaxUrl}/info`,
            form.serialize(),
            (data) => {
              if (data.Status === 0) {
                log('信息提交成功', data);
                // 在这里处理成功响应
              } else {
                log('信息提交失败', data);
                // 在这里处理错误响应
              }
            },
            'json'
          ).fail((xhr) => {
            log('信息提交请求失败', xhr.statusText);
          }).always(() => {
            // 重新启用表单控件
            if ($.utils && $.utils.disableFormInputs) {
              $.utils.disableFormInputs(form, false);
            }
          });
        },
  
        /**
         * 初始化插件
         */
        init() {
          log('初始化日记本插件');
          
          // 初始化表单
          methods.initDaybookForm();
          
          // 绑定所有事件
          settings.eventBindings.forEach(binding => {
            const handler = typeof binding.handler === 'string' ? 
              methods[binding.handler] : binding.handler;
              
            log(`绑定事件: ${binding.event} 到 ${binding.selector}`);
            
            $(document).on(binding.event, binding.selector, handler);
          });
          
          // 为HTML元素添加一个自定义属性，存储方法引用
          // 这样可以从DOM元素直接访问方法，而不需要全局对象
          this.data('diary', {
            doDaybook: methods.doDaybook,
            postInfo: methods.postInfo
          });
        }
      };
  
      // 对每个匹配的元素应用插件并添加方法访问
      return this.each(function() {
        log('应用到元素', this);
        methods.init.call($(this));
      });
    };
    
    // 添加公共方法直接到jQuery对象上，这样可以通过$.diary.方法名访问
    $.diary = {
      doDaybook: function() {
        const $el = $('body').data('diary');
        if ($el && $el.doDaybook) {
          $el.doDaybook.apply(this, arguments);
        }
      }
    };
  })(jQuery);