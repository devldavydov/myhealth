const template = `
<h3>${Constants.Page_Weight_WeightList}</h3>
<div class="row mb-2">
  <div class="col-sm-2">
    <a class="btn btn-primary" href="/weight/create" role="button"><i class="bi bi-plus-circle"></i></a>
    <button id="btnFilter" type="button" class="btn btn-primary"><i class="bi bi-funnel-fill"></i></button>
  </div>
</div>

${tmplLoader()}
${tmplToast()}

<div id="tblWeight" class="table-responsive d-none">
  <div class="mb-2">Диапазон: <span id="dateFrom"></span> - <span id="dateTo"></span></div>
  <table class="table table-striped table-bordered table-hover">
    <thead>
      <tr>
        <th class="align-middle col-2">${Constants.Common_Date}</th>
        <th class="align-middle col-9">${Constants.Weight_Value}</th>
        <th class="align-middle col-1 text-center"><i class="bi bi-gear"></i></th>
      </tr>
    </thead>
    <tbody>
    </tbody>
  </table>
</div>
`;

$( document ).ready(function() {
    createPage(template);

    let dateFrom = moment().subtract(1, 'year').valueOf();
    let dateTo = moment().valueOf();

    getWeightList(dateFrom, dateTo)
        .finally(() => {
            hideElement('#loader')
        })
        .then((data) => applyResult(data, dateFrom, dateTo))
        .catch((error) => {
            $('#toastBody').text(error.message);
            bootstrap.Toast.getOrCreateInstance($('#liveToast')).show();
        });
});

function applyResult(data, dateFrom, dateTo) {
    $('#dateFrom').text(moment(dateFrom).format('DD.MM.YYYY'));
    $('#dateTo').text(moment(dateTo).format('DD.MM.YYYY'));

    for (let w of data) {
        const tsStr = moment(w.timestamp).format('DD.MM.YYYY');
        $('#tblWeight tbody').append(`
            <tr>
                <td>${tsStr}</td>
                <td>${w.value}</td>
                <td class="align-middle text-center">
                    <a class="btn btn-sm btn-warning" href="edit?key=${w.timestamp}"><i class="bi bi-pencil"></i></a>
                </td>
            </tr>
        `);
    }

    showElement('#tblWeight');
}

async function getWeightList(from, to) {
    const resp = await axios.get(`api/list?from=${from}&to=${to}&order=desc`);
        
    if (resp.data.error) {
        throw new Error(resp.data.error);
    }

    let respWeight = [];
    for (let w of resp.data.data) {
        respWeight.push({
            timestamp: w.timestamp,
            value: w.value,
        });
    }

    return respWeight;
}