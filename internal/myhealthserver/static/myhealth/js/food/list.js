$( document ).ready(function() {
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
                <td>${f.cal100 }</td>
                <td class="myhealth-comment">${f.comment}</td>
                <td class="align-middle text-center"><a class="btn btn-sm btn-warning" href="edit/${f.key}"><i class="bi bi-pencil"></i></a></td>
            </tr>
        `);
    }
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