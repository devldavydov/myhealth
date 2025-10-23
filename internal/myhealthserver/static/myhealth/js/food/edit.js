const template = `
<h3>${Constants.Page_Food_FoodEdit}</h3>
${tmplLoader()}
${tmplToast()}
${tmplAlert('alert-danger d-none', 'alrt')}
<form id="frmEdit" class="d-none">
  <div class="mb-3">
    <label for="name" class="form-label">${Constants.Food_Name}</label>
    <input
      type="text"
      id="name"
      class="form-control"
      required>
  </div>
  <div class="mb-3">
    <label for="brand" class="form-label">${Constants.Food_Brand}</label>
    <input
      type="text"
      id="brand"
      class="form-control"
      required>
  </div>
  <div class="mb-3">
    <label for="cal100" class="form-label">${Constants.Food_Cal100}</label>
    <input
      type="number"
      id="cal100"
      min="0"
      step="0.01"
      class="form-control"
      required>
  </div>
  <div class="mb-3">
    <label for="prot100" class="form-label">${Constants.Food_Prot100}</label>
    <input
      type="number"
      id="prot100"
      min="0"
      step="0.01"
      class="form-control"
      required>
  </div>
  <div class="mb-3">
    <label for="fat100" class="form-label">${Constants.Food_Fat100}</label>
    <input
      type="number"
      id="fat100"
      min="0"
      step="0.01"
      class="form-control"
      required>
  </div>
  <div class="mb-3">
    <label for="carb100" class="form-label">${Constants.Food_Carb100}</label>
    <input
      type="number"
      id="carb100"
      min="0"
      step="0.01"
      class="form-control"
      required>
  </div>
  <div class="mb-3">
    <label for="comment" class="form-label">${Constants.Common_Comment}</label>
    <textarea
      type="text"
      id="comment"
      class="form-control"
      rows="3"
      required>
    </textarea>
  </div>
  <button type="submit" class="btn btn-primary mb-3">${Constants.Common_Save}</button>
</form>
`;

$( document ).ready(function() {
    createPage(template);

    const params = getQueryParams();
    if (!params.has('key')) {
        $('#alrt').text(Constants.Food_NotFound);
        showElement('#alrt');
        hideElement('#loader')
        return;
    }

    const foodKey = params.get('key');

    getFood(foodKey)
        .finally(() => {
            hideElement('#loader')
        })
        .then(setFood)
        .catch((error) => {
            $('#toastBody').text(error.message);
            bootstrap.Toast.getOrCreateInstance($('#liveToast')).show();
        });
});

async function getFood(key) {
    const resp = await axios.get(`api/get/${key}`);
        
    if (resp.data.error) {
        throw new Error(resp.data.error);
    }

    return resp.data;
}

function setFood(food) {
    $('#name').val(food.name);
    $('#brand').val(food.brand);
    $('#cal100').val(food.cal100);
    $('#prot100').val(food.prot100);
    $('#fat100').val(food.fat100);
    $('#carb100').val(food.carb100);
    $('#comment').val(food.comment);
    showElement('#frmEdit');
}