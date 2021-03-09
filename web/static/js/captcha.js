var onloadCallback = function() {
    grecaptcha.render('html_element', {
        'sitekey' : '6Le_BHgaAAAAAO36nUx-EwsaNYOvhJ90FtbHGCFb',
        'callback' : function(response) {
            const params = new URLSearchParams()
            params.append('g-recaptcha-response', response)
            const config = {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            }
            axios.post("/ok", params, config).then(response => {
                grecaptcha.reset(document.getElementById('html_element'));
                if (response.data.ok == true) {
                    document.getElementById("block").style.display = "none";
                } else {
                    document.getElementById("block").style.display = "block";
                }
            }).catch((err) => {
                alert(err);
            });
        },
        'theme' : 'light',
        'hl': 'ru'
    });
};
