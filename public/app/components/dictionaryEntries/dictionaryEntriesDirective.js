(function() {
    'use strict';
    angular
        .module('DVLF')
        .directive('dictionaryEntries', dictionaryEntries);

    function dictionaryEntries($http, $timeout, $log) {
        return {
            templateUrl: "app/components/dictionaryEntries/dictionaryEntries.html",
            link: function(scope, element, attrs) {
                var query = "/api/mot/" + attrs.head.trim();
                scope.dictionaries = [];
                var dicoOrder = ["feraud", "acad1694", "acad1762", "acad1798", "acad1835", "littre", "acad1932", "tlfi"];
                $http.get(query).then(function(response) {
                    for (var i = 0; i < dicoOrder.length; i++) {
                        if (dicoOrder[i] in response.data.dictionaries) {
                            scope.dictionaries.push({
                                "label": dicoOrder[i],
                                "data": response.data.dictionaries[dicoOrder[i]]
                            });
                        }
                    }
                });
            }
        }
    }
})();
