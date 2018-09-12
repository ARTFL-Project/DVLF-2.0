(function() {
    'use strict';
    angular
        .module('DVLF')
        .directive('nearestNeighbors', nearestNeighbors);

    function nearestNeighbors($rootScope, $timeout, $location) {
        return {
            templateUrl: "app/components/nearestNeighbors/nearestNeighbors.html",
            link: function(scope) {
                            
                scope.$on('resultsUpdate', function () {
                    $timeout(function() {
                        if (typeof(scope.nearestNeighbors) == 'undefined') { // avoid bubbling effect
                            scope.nearestNeighbors = scope.Results.results.nearestNeighbors;
                        }
                    });
                });
               
            }
        }
    }
})();