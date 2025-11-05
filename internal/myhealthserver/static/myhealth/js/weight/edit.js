const template = `
<h3>${Constants.Page_Weight_WeightEdit}</h3>
${tmplLoader()}
${tmplToast()}
${tmplAlert('alert-danger', 'alrt', '/weight')}
<form id="frmEdit" class="d-none">
  <div class="mb-3">
    <label for="timestamp" class="form-label">${Constants.Common_Date}</label>
    <input
      type="text"
      id="timestamp"
      class="form-control"
      disabled>
  </div>
  <div class="mb-3">
    <label for="value" class="form-label">${Constants.Weight_Value}</label>
    <input
      type="number"
      id="value"
      min="0.1"
      step="0.1"
      class="form-control"
      required>
  </div>  
  <div class="mb-3">
    <a href="/weight" class="btn btn-primary"><i class="bi bi-arrow-90deg-left"></i></a>
    <button type="submit" class="btn btn-primary"><i class="bi bi-floppy"></i></button>
    <button id="btnDelete" type="button" class="btn btn-danger"><i class="bi bi-trash"></i></button>
  </div>
</form>
`;

let weightKey = "";

$( document ).ready(function() {
    createPage(template);

    weightKey = (getQueryParams().get('key') || '').trim();
    if (weightKey === '' ) {
        $('#alrt .msg').text(Constants.Weight_NotFound);
        showElement('#alrt');
        hideElement('#loader')
        return;
    }

    getWeight(weightKey)
        .finally(() => {
            hideElement('#loader')
        })
        .then(setWeightForm)
        .catch((error) => {
            $('#alrt .msg').text(error.message);
            showElement('#alrt');
        });

    $('#frmEdit').submit(doEdit);
    $('#btnDelete').on('click', doDelete);
});

function setWeightForm(weight) {
    $('#timestamp').val(moment(weight.timestamp).format('DD.MM.YYYY'));
    $('#value').val(weight.value);

    showElement('#frmEdit');
}

function doEdit(e) {
    e.preventDefault();

    const value = parseFloat($('#value').val());

    setWeight({
        timestamp: parseInt(weightKey, 10),
        value: value
    })
        .then(() => {
            $('#toastBody').text(Constants.Weight_Saved);
            bootstrap.Toast.getOrCreateInstance($('#liveToast')).show();
        })
        .catch((error) => {
            $('#toastBody').text(error.message);
            bootstrap.Toast.getOrCreateInstance($('#liveToast')).show();
        });
}

function doDelete() {
    if (!confirm(Constants.Weight_DeletePrompt)) {
        return;
    }

    delWeight(weightKey)
        .then(() => {
            window.location.href = "/weight";
        })
        .catch((error) => {
            $('#toastBody').text(error.message);
            bootstrap.Toast.getOrCreateInstance($('#liveToast')).show();
        });
}

async function getWeight(key) {
    const resp = await axios.get(`api/${key}`);
        
    if (resp.data.error) {
        throw new Error(resp.data.error);
    }

    return resp.data;
}

async function setWeight(weight) {
    const resp = await axios.post('api/set', weight);
        
    if (resp.data.error) {
        throw new Error(resp.data.error);
    }

    return resp.data;
}

async function delWeight(key) {
    const resp = await axios.delete(`api/${key}`);
        
    if (resp.data.error) {
        throw new Error(resp.data.error);
    }

    return resp.data;
}