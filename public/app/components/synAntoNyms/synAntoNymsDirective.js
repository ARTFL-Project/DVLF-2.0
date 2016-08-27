(function() {
    'use strict';
    angular
        .module('DVLF')
        .directive('synAntoNyms', synAntoNyms);

    function synAntoNyms($rootScope, $timeout, $location) {
        return {
            templateUrl: "app/components/synAntoNyms/synAntoNyms.html",
            link: function(scope) {
                scope.$on('resultsUpdate', function () {
                    $timeout(function() {
                        if (typeof(scope.synonyms) == 'undefined') { // avoid bubbling effect
                            scope.synonyms = scope.Results.results.synonyms;
                            scope.antonyms = scope.Results.results.antonyms
                        }
                    });
                });
                scope.addSynonym = function() {
                    $location.url("/synonyme");
                }
                scope.addAntonym = function() {
                    $location.url("/antonyme");
                }
            }
        }
    }
})();
