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
                var dicoLabels = {
                    feraud: "Féraud: Dictionaire critique de la langue française (1787-1788)",
                    nicot: "Jean Nicot: Thresor de la langue française (1606)",
                    acad1694: "Dictionnaire de L'Académie française 1re édition (1694)",
                    acad1762: "Dictionnaire de L'Académie française 4e édition (1762)",
                    acad1798: "Dictionnaire de L'Académie française 5e édition (1798)",
                    acad1835: "Dictionnaire de L'Académie française 6e édition (1835)",
                    littre: "Émile Littré: Dictionnaire de la langue française (1872-1877)",
                    acad1932: "Dictionnaire de L'Académie française 8e édition (1932-1935)",
                    tlfi: "Le Trésor de la Langue Française Informatisé"
                }
                scope.dictionaries = [];
                var dicoOrder = ["tlfi", "acad1932", "littre", "acad1835", "acad1798", "feraud", "acad1762", "acad1694", "nicot"];
                $http.get(query).then(function(response) {
                    scope.Main.results = response.data;
                    for (var i = 0; i < dicoOrder.length; i++) {
                        if (dicoOrder[i] in response.data.dictionaries) {
                            scope.dictionaries.push({
                                "label": dicoLabels[dicoOrder[i]],
                                "data": response.data.dictionaries[dicoOrder[i]]
                            });
                        }
                    }
                });
            }
        }
    }
})();
