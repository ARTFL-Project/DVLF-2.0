(function() {
    'use strict';
    angular
        .module('DVLF')
        .directive('dictionaryEntries', dictionaryEntries);

    function dictionaryEntries($http, $timeout, $log, $rootScope) {
        var trimData = function(text) {
            var splitText = text.split(' ');
            var newText = "<span>" + splitText.slice(0, 50).join(' ') + "</span>";
            newText += '<span class="entry-end">' + splitText.slice(50, splitText.length).join(' ') + '</span>';
            return newText;
        }
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
                    tlfi: "Le Trésor de la Langue Française Informatisé",
                    bob: "BOB: Dictionaire d'argot"
                }
                scope.dictionaries = [];
                var dicoOrder = ["tlfi", "acad1932", "littre", "acad1835", "acad1798", "feraud", "acad1762", "acad1694", "nicot", "bob"];
                $http.get(query).then(function(response) {
                    scope.Main.results = response.data;
                    console.log(response)
                    $rootScope.$broadcast('resultsUpdate');
                    var displayed = 0;
                    var totalEntries = 0;
                    for (var i = 0; i < dicoOrder.length; i++) {
                        if (dicoOrder[i] in response.data.dictionaries) {
                            displayed += 1;
                            if (displayed < 3) {
                                var show = true;
                            } else {
                                var show = false;
                            }
                            totalEntries += response.data.dictionaries[dicoOrder[i]].length;
                            scope.dictionaries.push({
                                "name": dicoOrder[i],
                                "label": dicoLabels[dicoOrder[i]],
                                "data": response.data.dictionaries[dicoOrder[i]],
                                "show": show
                            });
                        }
                    }
                    if (scope.Main.results.userSubmit.length) {
                        totalEntries += scope.Main.results.userSubmit.length;
                        displayed += 1;
                        if (displayed < 3) {
                            var show = true;
                        } else {
                            var show = false;
                        }
                        scope.dictionaries.push({
                            name: "userSubmit",
                            label: "Définition(s) d'utilisateurs",
                            data: scope.Main.results.userSubmit,
                            show: show
                        })
                    }
                    scope.Results.totalDicos = displayed;
                    scope.Results.totalEntries = totalEntries;
                });
                scope.toggleEntry = function(event) {
                    var entry = angular.element(event.currentTarget).parents(".dico").find(".dico-entries");
                    var arrows = entry.parents(".dico").find(".arrow .glyphicon");
                    if (entry.css('display') === 'none') {
                        entry.addClass('show');
                        arrows.eq(1).hide();
                        arrows.eq(0).removeClass('ng-hide').show();
                    } else {
                        entry.removeClass("show");
                        arrows.eq(0).hide();
                        arrows.eq(1).removeClass('ng-hide').show();
                    }

                }
            }
        }
    }
})();
