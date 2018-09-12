(function() {
    'use strict';
    angular
        .module('DVLF')
        .directive('collocations', collocations);

    function collocations($rootScope, $timeout, $location) {
        return {
            templateUrl: "app/components/collocations/collocations.html",
            link: function(scope) {
                            
                scope.$on('resultsUpdate', function () {
                    $timeout(function() {
                        if (typeof(scope.collocations) == 'undefined') { // avoid bubbling effect
                            scope.collocations = scope.Results.results.collocates;
                        }
                    });
                });
               
            }
        }
    }
})();