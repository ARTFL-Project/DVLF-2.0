(function() {
    "use strict";

    angular
        .module("DVLF")
        .controller('NewSynAntoController', NewSynAntoController);

    function NewSynAntoController($scope, $log, $location, $routeParams, $http, $httpParamSerializer) {
        var vm = this;

        vm.typeOfNym = $location.path().replace('/', '');
        if (vm.typeOfNym == "synonyme") {
            var typeOfNym = "synonyms";
        } else {
            var typeOfNym = "antonyms";
        }
        vm.submission = {
            term: $scope.Main.queryTerm,
            nym: "",
            type: typeOfNym,
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
            $http.post("/api/submitNym", $httpParamSerializer(vm.submission), config)
                .then(
                    function(response) {
                        if (response.data.message === "success") {
                            $http.get('/api/wordwheel')
                                .then(
                                    function(response) {
                                        $scope.Main.wordwheel = response.data;
                                        $location.url("/mot/" + vm.submission.term)
                                    });
                        } else {
                            alert("Votre soumission a échouée.")
                        }
                    }
                );
        }

    }
})();
