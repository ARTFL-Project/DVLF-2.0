(function() {
    "use strict";

    angular
        .module("DVLF")
        .controller('ResultsController', ResultsController);

		function ResultsController($scope, $log, $location, $routeParams, totalResults) {
            var vm = this;
            vm.currentTerm = $routeParams.queryTerm;
            totalResults.queryTerm = vm.currentTerm;

            vm.define = function() {
                $location.path("/definition");
            }
        }
})();
