$(() => {
  const AJAX_URL = '';

  $('body').passport({ ajaxUrl: AJAX_URL });
  $('body').diary({ ajaxUrl: AJAX_URL });

  // Header fixed
  $(document).scroll(() => {
    if ($(document).scrollTop() > 1) {
      $('#top-nav').html($('#page-nav').html());
    } else if ($(document).scrollTop() === 0) {
      $('#top-nav').html('');
    }
  });
});
