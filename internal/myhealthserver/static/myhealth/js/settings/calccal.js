const template = `
<h3>${Constants.Page_Settings_CalcCal}</h3>
<form id="calc">
  <div class="mb-3">
    <label class="form-label">${Constants.Settings_Gender}</label>
    <div class="form-check">
      <input
        class="form-check-input"
        name="gender"
        value="m"
        type="radio"
        id="genderM"
        checked
      >
      <label class="form-check-label" for="genderM">${Constants.Settings_GenderM}</label>
    </div>
    <div class="form-check">
      <input
        class="form-check-input"
        name="gender"
        value="f"
        id="genderF"
        type="radio"
      >
      <label class="form-check-label" for="genderF">${Constants.Settings_GenderF}</label>
    </div>
  </div>
  <div class="mb-3">
    <label for="weight" class="form-label">${Constants.Settings_Weight}</label>
    <input
      type="number"
      id="weight"
      min="0.1"
      step="0.1"
      class="form-control"
      required>
  </div>
  <div class="mb-3">
    <label for="height" class="form-label">${Constants.Settings_Height}</label>
    <input
      type="number"
      id="height"
      min="0.1"
      step="0.1"
      class="form-control"
      required>
  </div>
  <div class="mb-3">
    <label for="age" class="form-label">${Constants.Settings_Age}</label>
    <input
      type="number"
      id="age"
      min="1"
      step="1"
      class="form-control"
      required>
  </div>
  <button type="submit" class="btn btn-primary">${Constants.Common_Calc}</button>
</form>

<table id="tblResult" class="table mt-3 d-none">
  <tbody>
  </tbody>
</table>
`;

$( document ).ready(function() {
    createPage(template);

    $('#calc').submit(doCalc);
});

function doCalc(e) {
    e.preventDefault();

    const gender = $('input[name="gender"]:checked').val();
    const weight = parseFloat($('#weight').val());
    const height = parseFloat($('#height').val());
    const age = parseInt($('#age').val(), 10);
   
    let ubm = 10 * weight + 6.25 * height - 5 * age;
    if (gender === "m") 
        ubm += 5;
    else
        ubm -= 161;

    const activities = [
        {name: Constants.Settings_Ubm, k: 1},
        {name: Constants.Settings_Activity1, k: 1.2},
        {name: Constants.Settings_Activity2, k: 1.375},
        {name: Constants.Settings_Activity3, k: 1.55},
        {name: Constants.Settings_Activity4, k: 1.725},
        {name: Constants.Settings_Activity5, k: 1.9}
    ];

    $('#tblResult tbody').empty();
    activities.forEach(el => {
        $('#tblResult tbody').append(`<tr><th>${el.name}</th><td>${(ubm * el.k).toFixed(2)}</td>`)
    });
    showElement('#tblResult');
}