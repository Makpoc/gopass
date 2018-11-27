function copyToCB() {
  var password = document.querySelector('.password');
  var passwordBtn = document.querySelector('.copyToClipboard');
  var selection = window.getSelection();
  var range = document.createRange();
  range.selectNodeContents(password);
  selection.removeAllRanges();
  selection.addRange(range);

  try {
    document.execCommand('copy');
    selection.removeAllRanges();
    passwordBtn.textContent = '✓ Copied to clipboard';
  } catch(e) {
    passwordBtn.textContent = '✗ Failed to copy to clipboard';
    console.log(e)
  }
};
