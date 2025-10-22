const template = `
<div class="alert alert-danger" role="alert">
  ${Constants.Page_NotFound}
</div>
`;

$( document ).ready(function() {
    createPage(template);
});
