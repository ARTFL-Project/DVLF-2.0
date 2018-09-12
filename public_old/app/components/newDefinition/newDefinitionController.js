(function() {
    "use strict";

    angular
        .module("DVLF")
        .controller('NewDefinitionController', NewDefinitionController);

    function NewDefinitionController($scope, $log, $location, $routeParams, $http, $httpParamSerializer, totalResults) {
        var vm = this;

        vm.submission = {
            term: totalResults.queryTerm,
            definition: "",
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
            $http.post("/api/submit", $httpParamSerializer(vm.submission), config)
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
