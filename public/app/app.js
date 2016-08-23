(function() {
    "use strict";

    angular
        .module("DVLF", ['ngRoute', 'ngTouch', 'ngSanitize', 'ngAnimate', 'angucomplete-alt', 'sticky', 'vcRecaptcha']);

    angular
        .module("DVLF")
        .controller('MainController', MainController);

		function MainController($scope, $routeParams, $route, $log, $http, $location) {
            var vm = this;

            if($routeParams.queryTerm !== 'undefined') {
                vm.queryTerm = $routeParams.queryTerm;
            } else {
                vm.queryTerm = "";
            }


            $http.get('/api/wordwheel').then(function(response) {
                vm.wordwheel = response.data;
                vm.wordwheelObject = {};
                for (var i=0; i < vm.wordwheel.length; i+=1) {
                    vm.wordwheelObject[vm.wordwheel[i]] = i;
                }
            });

            vm.apropos = false;
            vm.viewAPropos = function() {
                $location.path("/apropos")
            }

            vm.search = function(word) {
                // if (typeof(word) === 'undefined') {
                    word = angular.element("#search_value").val();
                    $location.path("/mot/" + word.trim());
                //     console.log(word)
                //     $scope.$apply();
                // } else {
                //     console.log(word)
                //     $location.path("/mot/" + word.trim());
                // }
            }

            $scope.autocomplete = function(query) {
                if (typeof(query) !== 'undefined') {
                    $location.path("/mot/" + query.title.trim());
                }
            }
        }
})();
