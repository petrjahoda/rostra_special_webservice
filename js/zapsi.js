const userRow = document.getElementById("user-row");
const orderRow = document.getElementById("order-row");
const operationRow = document.getElementById("operation-row");
const workplaceRow = document.getElementById("workplace-row");


const resetButton = document.getElementById("reset-button");

const userInput = document.getElementById("user-input");
const orderInput = document.getElementById("order-input");
const countOkInput = document.getElementById("count-ok-input");
const countNokInput = document.getElementById("count-nok-input");
const operationSelect = document.getElementById("operation-select");
const workplaceSelect = document.getElementById("workplace-select");
const nokTypeSelect = document.getElementById("nok-type-select");

const userOkButton = document.getElementById("user-ok-button");
const orderOkButton = document.getElementById("order-ok-button");
const operationOkButton = document.getElementById("operation-ok-button");
const workplaceOkButton = document.getElementById("workplace-ok-button");
const countButton = document.getElementById("count-button");

const orderBackButton = document.getElementById("order-back-button");
const operationBackButton = document.getElementById("operation-back-button");
const workplaceBackButton = document.getElementById("workplace-back-button");
const countBackButton = document.getElementById("count-back-button");

const clovekRadio = document.getElementById("clovek-radio");
const serizeniRadio = document.getElementById("serizeni-radio");
const strojRadio = document.getElementById("stroj-radio");

const startOrderButton = document.getElementById("start-order-button");
const endOrderButton = document.getElementById("end-order-button");
const transferOrderButton = document.getElementById("transfer-order-button");

const table = document.getElementById("table");

const errorInfoPanel = document.getElementById("error-info-panel");

resetButton.addEventListener("click", () => {
    window.location.replace('');
})

sessionStorage.clear()


