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
                scope.dictionaries = [];
                scope.Main.loading = true;
                $http.get(query).then(function(response) {
                    scope.Results.results = response.data;
                    scope.totalEntries = response.data.dictionaries.totalEntries;
                    $rootScope.$broadcast('resultsUpdate', {
                        totalDicos: response.data.dictionaries.totalDicos,
                        totalEntries: response.data.dictionaries.totalEntries,
                        queryTerm: attrs.head.trim()
                    });
                    scope.Main.loading = false;
                    var displayed = 0;
                    var totalEntries = 0;
                    if (response.data != null) {

                    } else {
                        scope.Results.results = {
                            totalEntries: 0
                        };
                    }
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
