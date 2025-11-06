const template = `
<h3>${Constants.Page_Food_FoodEdit}</h3>
${tmplLoader()}
${tmplToast()}
${tmplAlert('alert-danger', 'alrt', '/food')}
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
      class="form-control">
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
      rows="3">
    </textarea>
  </div>
  <div class="mb-3">
    <a href="/food" class="btn btn-primary"><i class="bi bi-arrow-90deg-left"></i></a>
    <button type="submit" class="btn btn-primary"><i class="bi bi-floppy"></i></button>
    <button id="btnDelete" type="button" class="btn btn-danger"><i class="bi bi-trash"></i></button>
  </div>
</form>
`;

let foodKey = "";

$( document ).ready(function() {
    createPage(template);

    foodKey = (getQueryParams().get('key') || '').trim();
    if (foodKey === '' ) {
        showAlert('alrt', Constants.Food_NotFound);
        hideElement('#loader');
        return;
    }
    
    getFood(foodKey)
        .finally(() => {
            hideElement('#loader')
        })
        .then(setFoodForm)
        .catch((error) => {
            showAlert('alrt', error.message);
        });

    $('#frmEdit').submit(doEdit);
    $('#btnDelete').on('click', doDelete);
});

function setFoodForm(food) {
    $('#name').val(food.name);
    $('#brand').val(food.brand);
    $('#cal100').val(food.cal100);
    $('#prot100').val(food.prot100);
    $('#fat100').val(food.fat100);
    $('#carb100').val(food.carb100);
    $('#comment').val(food.comment);

    showElement('#frmEdit');    
}

function doEdit(e) {
    e.preventDefault();

    const name = $('#name').val().trim();
    const brand = $('#brand').val().trim();
    const cal100 = parseFloat($('#cal100').val());
    const prot100 = parseFloat($('#prot100').val());
    const fat100 = parseFloat($('#fat100').val());
    const carb100 = parseFloat($('#carb100').val());
    const comment = $('#comment').val().trim();

    setFood({
        key: foodKey,
        name: name,
        brand: brand,
        cal100: cal100,
        prot100: prot100,
        fat100: fat100,
        carb100: carb100,
        comment: comment
    })
        .then(() => {
            $('#toastBody').text(Constants.Food_Saved);
            bootstrap.Toast.getOrCreateInstance($('#liveToast')).show();
        })
        .catch((error) => {
            $('#toastBody').text(error.message);
            bootstrap.Toast.getOrCreateInstance($('#liveToast')).show();
        });
}

function doDelete() {
    if (!confirm(Constants.Food_DeletePrompt)) {
        return;
    }

    delFood(foodKey)
        .then(() => {
            window.location.href = "/food";
        })
        .catch((error) => {
            $('#toastBody').text(error.message);
            bootstrap.Toast.getOrCreateInstance($('#liveToast')).show();
        });
}

async function getFood(key) {
    const resp = await axios.get(`api/${key}`);
        
    if (resp.data.error) {
        throw new Error(resp.data.error);
    }

    return resp.data;
}

async function setFood(food) {
    const resp = await axios.post('api/set', food);
        
    if (resp.data.error) {
        throw new Error(resp.data.error);
    }

    return resp.data;
}

async function delFood(key) {
    const resp = await axios.delete(`api/${key}`);
        
    if (resp.data.error) {
        throw new Error(resp.data.error);
    }

    return resp.data;
}