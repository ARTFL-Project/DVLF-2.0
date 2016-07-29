(function() {
    'use strict';
    angular
        .module('DVLF')
        .directive('synAntoNyms', synAntoNyms);

    function synAntoNyms($rootScope, $timeout) {
        var buildSynAntoLists = function(scope) {
            scope.synonyms = [];
            for (var i=0; i < scope.Main.results.synonyms.length; i++) {
                var link = '<a href="/mot/' + scope.Main.results.synonyms[i] + '">' + scope.Main.results.synonyms[i] + '</a>';
                scope.synonyms.push(link)
            }
            scope.synonyms = scope.synonyms.join(', ');
            scope.antonyms = [];
            for (var i=0; i < scope.Main.results.antonyms.length; i++) {
                var link = '<a href="/mot/' + scope.Main.results.antonyms[i] + '">' + scope.Main.results.antonyms[i] + '</a>';
                scope.antonyms.push(link)
            }
            scope.antonyms = scope.antonyms.join(', ');
        }
        return {
            templateUrl: "app/components/synAntoNyms/synAntoNyms.html",
            link: function(scope) {
                buildSynAntoLists(scope);
                scope.$on('resultsUpdate', function () {
                    $timeout(function() {
                        buildSynAntoLists(scope);
                    });
                });
            }
        }
    }
})();
