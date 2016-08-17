(function() {
    "use strict";

    angular
        .module("DVLF")
        .controller('NewExampleController', NewExampleController);

    function NewExampleController($scope, $log, $location, $routeParams, $http, $httpParamSerializer) {
        var vm = this;

        vm.submission = {
            term: $scope.Main.queryTerm,
            example: "",
            source: "",
            link: "",
            recaptchaResponse: ""
        };

        vm.recaptcha = {
            key: "6LfhfycTAAAAAId87HWFIW-N8cShdp6O8fpAMK8h"
        }

        vm.submitting = false;

        vm.recaptchaResponse = function(response) {
            vm.submission.recaptchaResponse = response;
        }

        vm.submit = function() {
            vm.submitting = true;
            var config = {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded;charset=utf-8;'
                }
            }
            $http.post("/api/submitExample", $httpParamSerializer(vm.submission), config)
                .then(
                    function(response) {
                        if (response.data.message === "success") {
                            $location.url("/mot/" + vm.submission.term)
                        } else {
                            alert("Votre soumission a échouée.")
                        }
                    }
                );
        }

    }
})();
