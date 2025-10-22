const template = `
<h3>${Constants.Page_Finance_BondCalc}</h3>
<form id="calc">
  <div class="mb-3">
    <label for="nominal" class="form-label">${Constants.Finance_Nominal}</label>
    <input
      type="number"
      id="nominal"
      min="0"
      step="0.01"
      class="form-control"
      value="1000"
      required>
  </div>
  <div class="mb-3">
    <label for="price" class="form-label">${Constants.Finance_BondCurrentPrice}</label>
    <input
      type="number"
      id="price"
      min="0"
      step="0.01"
      class="form-control"
      required>
  </div>
  <div class="mb-3">
    <label for="matDate" class="form-label">${Constants.Finance_BondMaturityDate}</label>
    <input
      type="date"
      id="matDate"
      class="form-control"
      required>
    <div id="matDateHelp" class="form-text text-danger d-none">${Constants.Common_DateGToday}</div>
  </div>
  <div class="mb-3">
    <label for="ncd" class="form-label">${Constants.Finance_BondNcd}</label>
    <input
      type="number"
      id="ncd"
      min="0"
      step="0.01"
      class="form-control"
      required>
  </div>
  <div class="mb-3">
    <label for="cd" class="form-label">${Constants.Finance_BondCd}</label>
    <input
      type="number"
      id="cd"
      min="0"
      step="0.01"
      class="form-control"
      required>
  </div>
  <div class="mb-3">
    <label for="cdCnt" class="form-label">${Constants.Finance_BondCdCount}</label>
    <input
      type="number"
      id="cdCnt"
      min="0"
      step="1"
      class="form-control"
      required>
  </div>
  <div class="mb-3">
    <label for="sum" class="form-label">${Constants.Finance_BondSum}</label>
    <input
      type="number"
      id="sum"
      min="0"
      step="0.01"
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

    const nominal = parseFloat($('#nominal').val());
    const price = parseFloat($('#price').val());
    const matDate = new Date($('#matDate').val());
    const ncd = parseFloat($('#ncd').val());
    const cd = parseFloat($('#cd').val());
    const cdCnt = parseFloat($('#cdCnt').val());
    const sum = parseFloat($('#sum').val());

    if (!checkForm(matDate)) {
        return;
    }

    const totalPrice = price + ncd;
    const cdSum = cd * cdCnt;
    const totalBoughtCnt = Math.round(sum / totalPrice);
    const totalBoughtSum = totalBoughtCnt * totalPrice;
    const daysToMaturity = getDaysBetween(matDate, new Date());
    const ytm = ((nominal + cdSum - totalPrice) / totalPrice) * (365 / daysToMaturity) * 100;
    const totalSum = totalBoughtCnt * (nominal + cdSum);
    const diffSum = totalSum - totalBoughtSum;

    const result = [
        {name: Constants.Finance_TotalBoughtCnt, val: totalBoughtCnt.toString()},
        {name: Constants.Finance_TotalBoughtSum, val: totalBoughtSum.toFixed(2)},
        {name: Constants.Finance_DaysToMaturity, val: daysToMaturity.toString()},
        {name: Constants.Finance_YTM, val: ytm.toFixed(2)},
        {name: Constants.Finance_TotalSum, val: totalSum.toFixed(2)},
        {name: Constants.Finance_DiffSum, val: diffSum.toFixed(2)},
    ];

    $('#tblResult tbody').empty();
    result.forEach(el => {
        $('#tblResult tbody').append(`<tr><th>${el.name}</th><td>${el.val}</td>`)
    });
    showElement('#tblResult');
}

function checkForm(matDate) {
    let res = true;

    if (new Date() > matDate) {
      showElement('#matDateHelp');
      res = false;
    } else {
      hideElement('#matDateHelp');
    }

    return res;
}

function getDaysBetween(startDate, endDate) {
    const date1 = new Date(startDate);
    const date2 = new Date(endDate);

    const diffInMilliseconds = Math.abs(date2.getTime() - date1.getTime());
    const oneDayInMilliseconds = 1000 * 60 * 60 * 24;

    const diffInDays = Math.round(diffInMilliseconds / oneDayInMilliseconds);
    return diffInDays;
}