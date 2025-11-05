const template = `
<h3>${Constants.Page_Weight_WeightCreate}</h3>
${tmplToast()}
<form id="frmCreate">
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
  </div>
</form>
`;

let weightKey = moment().startOf('day').valueOf();

$( document ).ready(function() {
    createPage(template);

    $('#timestamp').val(moment(weightKey).format('DD.MM.YYYY'));
    $('#frmCreate').submit(doCreate);
});

function doCreate(e) {
    e.preventDefault();

    const value = parseFloat($('#value').val());

    setWeight({
        timestamp: weightKey,
        value: value
    })
        .then(() => {
            window.location.href = "/weight";
        })
        .catch((error) => {
            $('#toastBody').text(error.message);
            bootstrap.Toast.getOrCreateInstance($('#liveToast')).show();
        });
}

async function setWeight(weight) {
    const resp = await axios.post('api/set', weight);
        
    if (resp.data.error) {
        throw new Error(resp.data.error);
    }

    return resp.data;
}