$( document ).ready(function() {
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
        {name: Constants['Settings_Ubm'], k: 1},
        {name: Constants['Settings_Activity1'], k: 1.2},
        {name: Constants['Settings_Activity2'], k: 1.375},
        {name: Constants['Settings_Activity3'], k: 1.55},
        {name: Constants['Settings_Activity4'], k: 1.725},
        {name: Constants['Settings_Activity5'], k: 1.9}
    ];

    $('#tblResult tbody').empty();
    activities.forEach(el => {
        $('#tblResult tbody').append(`<tr><th>${el.name}</th><td>${(ubm * el.k).toFixed(2)}</td>`)
    });
    showElement('#tblResult');
}