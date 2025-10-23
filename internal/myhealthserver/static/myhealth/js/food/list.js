const template = `
<h3>${Constants.Page_Food_FoodList}</h3>
<div class="row mb-2">
  <div class="col-sm-2">
    <a class="btn btn-primary" href="/food/create" role="button"><i class="bi-plus-square"></i> ${Constants.Common_Create}</a>  
  </div>
</div>
<div class="row mb-2">
${tmplSearch()}
</div>

${tmplLoader()}
${tmplToast()}

<div id="tblFood" class="table-responsive d-none">
  <table class="table table-striped table-bordered table-hover">
    <thead>
      <tr>
        <th class="align-middle col-4">${Constants.Food_Name}</th>
        <th class="align-middle col-2">${Constants.Food_Brand}</th>
        <th class="align-middle col-1">${Constants.Food_CPFC}</th>
        <th class="align-middle col-2">${Constants.Food_Comment}</th>
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

    getFoodList()
        .finally(() => {
            hideElement('#loader')
        })
        .then(applyResult)
        .catch((error) => {
            $('#toastBody').text(error.message);
            bootstrap.Toast.getOrCreateInstance($('#liveToast')).show();
        });

    $("#search").on("keyup", search);
    $("#btnSearchClear").on('click', function() {
        $('#search').val('');
        $('#search').trigger('keyup');
    });
});

async function getFoodList() {
    const resp = await axios.get("api/list");
        
    if (resp.data.error) {
        throw new Error(resp.data.error);
    }

    let respFood = [];
    for (let f of resp.data.data) {
        respFood.push({
            key: f.key,
            name: f.name,
            brand: f.brand,
            cal100: f.cal100,
            prot100: f.prot100,
            fat100: f.fat100,
            carb100: f.carb100,
            comment: f.comment,
        });
    }

    return respFood;
}

function applyResult(data) {
    for (let f of data) {
        $('#tblFood tbody').append(`
            <tr>
                <td class="myhealth-name">${f.name}</td>
                <td class="myhealth-brand">${f.brand }</td>
                <td class="align-middle text-center">
                    <button
                        class="btn btn-sm btn-info"
                        data-bs-toggle="popover"
                        data-bs-content="<strong>${Constants.Food_C}</strong>: ${f.cal100}<br/><strong>${Constants.Food_P}</strong>: ${f.prot100}<br/><strong>${Constants.Food_F}</strong>: ${f.fat100}<br/><strong>${Constants.Food_Cb}</strong>: ${f.carb100}"
                        data-bs-html="true"
                    >
                        <i class="bi bi-info-circle"></i>
                    </button>
                </td>
                <td class="myhealth-comment">${f.comment}</td>
                <td class="align-middle text-center"><a class="btn btn-sm btn-warning" href="edit/${f.key}"><i class="bi bi-pencil"></i></a></td>
            </tr>
        `);
    }

    $('[data-bs-toggle="popover"]').each((_, el) => {
        new bootstrap.Popover(el);
    });

    showElement('#tblFood');
}

function search() {
    const pattern = $(this).val().toLocaleUpperCase();

    $("#tblFood tr").each(function(index) {
        if (index === 0) {
            return;
        }

        $row = $(this);

        let name = $row.find("td.myhealth-name").text().toLocaleUpperCase();
        let brand = $row.find("td.myhealth-brand").text().toLocaleUpperCase();
        let comment = $row.find("td.myhealth-comment").text().toLocaleUpperCase();

        if (name.indexOf(pattern) === -1 && 
            brand.indexOf(pattern) === -1 &&
            comment.indexOf(pattern) === -1) {
            $row.hide();
        }
        else {
            $row.show();
        }
    });
}